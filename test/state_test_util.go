package test

import (
	"bytes"
	"fmt"
	"io"
	"math/big"
	"strconv"

	"github.com/MonteCarloClub/KBD/block_error"
	"github.com/MonteCarloClub/KBD/common"
	"github.com/MonteCarloClub/KBD/common/logger/glog"
	"github.com/MonteCarloClub/KBD/crypto"
	"github.com/MonteCarloClub/KBD/kbpool"
	"github.com/MonteCarloClub/KBD/kdb"
	"github.com/MonteCarloClub/KBD/state"
	"github.com/MonteCarloClub/KBD/vm"
)

func RunStateTestWithReader(r io.Reader, skipTests []string) error {
	tests := make(map[string]VmTest)
	if err := readJson(r, &tests); err != nil {
		return err
	}

	if err := runStateTests(tests, skipTests); err != nil {
		return err
	}

	return nil
}

func RunStateTest(p string, skipTests []string) error {
	tests := make(map[string]VmTest)
	if err := readJsonFile(p, &tests); err != nil {
		return err
	}

	if err := runStateTests(tests, skipTests); err != nil {
		return err
	}

	return nil

}

func runStateTests(tests map[string]VmTest, skipTests []string) error {
	skipTest := make(map[string]bool, len(skipTests))
	for _, name := range skipTests {
		skipTest[name] = true
	}

	for name, test := range tests {
		if skipTest[name] {
			glog.Infoln("Skipping state test", name)
			return nil
		}

		if err := runStateTest(test); err != nil {
			return fmt.Errorf("%s: %s\n", name, err.Error())
		}

		glog.Infoln("State test passed: ", name)
		//fmt.Println(string(statedb.Dump()))
	}
	return nil

}

func runStateTest(t VmTest) error {
	db, _ := kdb.NewMemDatabase()
	statedb := state.New(common.Hash{}, db)
	for addr, account := range t.Pre {
		obj := StateObjectFromAccount(db, addr, account)
		statedb.SetStateObject(obj)
		for a, v := range account.Storage {
			obj.SetState(common.HexToHash(a), common.HexToHash(v))
		}
	}

	// XXX Yeah, yeah...
	env := make(map[string]string)
	env["currentCoinbase"] = t.Env.CurrentCoinbase
	env["currentDifficulty"] = t.Env.CurrentDifficulty
	env["currentGasLimit"] = t.Env.CurrentGasLimit
	env["currentNumber"] = t.Env.CurrentNumber
	env["previousHash"] = t.Env.PreviousHash
	if n, ok := t.Env.CurrentTimestamp.(float64); ok {
		env["currentTimestamp"] = strconv.Itoa(int(n))
	} else {
		env["currentTimestamp"] = t.Env.CurrentTimestamp.(string)
	}

	var (
		ret []byte
		// gas  *big.Int
		// err  error
		logs state.Logs
	)

	ret, logs, _, _ = RunState(statedb, env, t.Transaction)

	// // Compare expected  and actual return
	rexp := common.FromHex(t.Out)
	if bytes.Compare(rexp, ret) != 0 {
		return fmt.Errorf("return failed. Expected %x, got %x\n", rexp, ret)
	}

	// check post state
	for addr, account := range t.Post {
		obj := statedb.GetStateObject(common.HexToAddress(addr))
		if obj == nil {
			continue
		}

		if obj.Balance().Cmp(common.Big(account.Balance)) != 0 {
			return fmt.Errorf("(%x) balance failed. Expected %v, got %v => %v\n", obj.Address().Bytes()[:4], account.Balance, obj.Balance(), new(big.Int).Sub(common.Big(account.Balance), obj.Balance()))
		}

		if obj.Nonce() != common.String2Big(account.Nonce).Uint64() {
			return fmt.Errorf("(%x) nonce failed. Expected %v, got %v\n", obj.Address().Bytes()[:4], account.Nonce, obj.Nonce())
		}

		for addr, value := range account.Storage {
			v := obj.GetState(common.HexToHash(addr))
			vexp := common.HexToHash(value)

			if v != vexp {
				return fmt.Errorf("(%x: %s) storage failed. Expected %x, got %x (%v %v)\n", obj.Address().Bytes()[0:4], addr, vexp, v, vexp.Big(), v.Big())
			}
		}
	}

	statedb.Sync()
	if common.HexToHash(t.PostStateRoot) != statedb.Root() {
		return fmt.Errorf("Post state root error. Expected %s, got %x", t.PostStateRoot, statedb.Root())
	}

	// check logs
	if len(t.Logs) > 0 {
		if err := checkLogs(t.Logs, logs); err != nil {
			return err
		}
	}

	return nil
}

func RunState(statedb *state.StateDB, env, tx map[string]string) ([]byte, state.Logs, *big.Int, error) {
	reader := bytes.NewReader([]byte(tx["secretKey"]))
	var (
		keyPair = crypto.NewKey(reader)
		data    = common.FromHex(tx["data"])
		gas     = common.Big(tx["gasLimit"])
		price   = common.Big(tx["gasPrice"])
		value   = common.Big(tx["value"])
		nonce   = common.Big(tx["nonce"]).Uint64()
		caddr   = common.HexToAddress(env["currentCoinbase"])
	)

	var to *common.Address
	if len(tx["to"]) > 2 {
		t := common.HexToAddress(tx["to"])
		to = &t
	}
	// Set pre compiled contracts
	vm.Precompiled = vm.PrecompiledContracts()

	snapshot := statedb.Copy()
	coinbase := statedb.GetOrNewStateObject(caddr)
	coinbase.SetGasLimit(common.Big(env["currentGasLimit"]))

	message := NewMessage(keyPair.Address, to, data, value, gas, price, nonce)
	vmenv := NewEnvFromMap(statedb, env, tx)
	vmenv.origin = keyPair.Address
	ret, _, err := kbpool.ApplyMessage(vmenv, message, coinbase)
	if block_error.IsNonceErr(err) || block_error.IsInvalidTxErr(err) || state.IsGasLimitErr(err) {
		statedb.Set(snapshot)
	}
	statedb.Sync()

	return ret, vmenv.state.Logs(), vmenv.Gas, err
}
