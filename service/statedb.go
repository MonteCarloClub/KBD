package service

import (
	"context"
	"github.com/MonteCarloClub/KBD/common"
	"github.com/MonteCarloClub/KBD/constant"
	"github.com/MonteCarloClub/KBD/kdb"
	"github.com/MonteCarloClub/KBD/state"
	"path"
)

func GetAccountData(ctx context.Context, address string) *state.StateObject {
	file := path.Join("/", constant.DBDir, constant.DBFile)
	db, _ := kdb.NewLDBDatabase(file)
	statedb := state.New(common.Hash{}, db)
	obj := statedb.GetStateObject(common.StringToAddress(address))
	db.Close()
	return obj
}
