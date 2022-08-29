package service

import (
	"context"

	"github.com/MonteCarloClub/KBD/model/state"

	"github.com/MonteCarloClub/KBD/common"
	"github.com/MonteCarloClub/KBD/frame"
)

func GetAccountData(ctx context.Context, address string) *state.StateObject {
	stateDB := frame.GetStateDB(ctx)
	obj := stateDB.GetStateObject(common.StringToAddress(address))
	return obj
}
