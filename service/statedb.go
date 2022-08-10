package service

import (
	"context"
	"github.com/MonteCarloClub/KBD/common"
	"github.com/MonteCarloClub/KBD/kdb"
	"github.com/MonteCarloClub/KBD/state"
)

func GetAccountData(ctx context.Context, address string) *state.StateObject {
	db, _ := kdb.NewMemDatabase()
	statedb := state.New(common.Hash{}, db)
	obj := statedb.GetStateObject(common.StringToAddress(address))
	return obj
}
