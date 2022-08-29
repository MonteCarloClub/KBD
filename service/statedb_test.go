package service

import (
	"fmt"
	"testing"

	"github.com/MonteCarloClub/KBD/common"
	"github.com/MonteCarloClub/KBD/frame"
	"github.com/MonteCarloClub/KBD/model/state"
)

func TestPutAccountData(t *testing.T) {
	stateDB := frame.GetStateDB()
	address := common.StringToAddress("0x945304eb96065b2a98b57a48a06ae28d285a71b5")
	account := Account{
		Balance: "0x17",
		Code:    "0x6000355415600957005b60203560003555",
		Nonce:   "0x00",
		Storage: nil,
	}
	obj := StateObjectFromAccount(frame.GetDB(), address, account)
	stateDB.UpdateStateObject(obj)
	stateDB.Trie().Commit()
	res := stateDB.GetStateObject(address)
	fmt.Println(res)
}

func TestGetAccountData(t *testing.T) {
	stateDB := frame.GetStateDB()
	address := common.StringToAddress("0x945304eb96065b2a98b57a48a06ae28d285a71b5")
	res := stateDB.GetStateObject(address)
	fmt.Println(res)
}

type Account struct {
	Balance string
	Code    string
	Nonce   string
	Storage map[string]string
}

func StateObjectFromAccount(db common.Database, address common.Address, account Account) *state.StateObject {
	obj := state.NewStateObject(address, db)
	obj.SetBalance(common.Big(account.Balance))

	if common.IsHex(account.Code) {
		account.Code = account.Code[2:]
	}
	obj.SetCode(common.Hex2Bytes(account.Code))
	obj.SetNonce(common.Big(account.Nonce).Uint64())

	return obj
}
