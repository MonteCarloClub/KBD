package service

import (
	"context"
	"path"

	"github.com/MonteCarloClub/KBD/common"
	"github.com/MonteCarloClub/KBD/constant"
	"github.com/MonteCarloClub/KBD/kdb"
	"github.com/MonteCarloClub/KBD/state"
)

func GetAccountData(ctx context.Context, address string) *state.StateObject {
	file := path.Join("/", constant.DataDir, constant.StateDBFile)
	db, _ := kdb.NewLDBDatabase(file)
	statedb := state.New(common.Hash{}, db)
	obj := statedb.GetStateObject(common.StringToAddress(address))
	db.Close()
	return obj
}
