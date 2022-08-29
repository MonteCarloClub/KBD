package service

import (
	"context"

	"github.com/MonteCarloClub/KBD/kitex_gen/api"

	"github.com/MonteCarloClub/KBD/model/state"

	"github.com/MonteCarloClub/KBD/common"
	"github.com/MonteCarloClub/KBD/frame"
)

func SetAccountData(ctx context.Context, req *api.SetAccountDataRequest) bool {
	stateDB := frame.GetStateDB()
	address := common.HexToAddress(req.Address)
	obj := StateObjectFromAccount(frame.GetDB(), address, req.GetBalance(), req.GetCode(), req.GetNonce())
	stateDB.Trie().Commit()
	stateDB.UpdateStateObject(obj)
	return true
}

func GetAccountData(ctx context.Context, address string) *state.StateObject {
	stateDB := frame.GetStateDB()
	obj := stateDB.GetStateObject(common.HexToAddress(address))
	return obj
}

func StateObjectFromAccount(db common.Database, address common.Address, balance string, code string, nonce string) *state.StateObject {
	obj := state.NewStateObject(address, db)
	obj.SetBalance(common.Big(balance))

	if common.IsHex(code) {
		code = code[2:]
	}
	obj.SetCode(common.Hex2Bytes(code))
	obj.SetNonce(common.Big(nonce).Uint64())

	return obj
}
