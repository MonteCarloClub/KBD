package model

import (
	"github.com/MonteCarloClub/KBD/kitex_gen/api"
	"github.com/MonteCarloClub/KBD/model/state"
)

func StateObject2VO(obj *state.StateObject) *api.Account {
	if obj == nil {
		return nil
	}
	return &api.Account{
		Address: obj.Address().Str(),
		Balance: obj.Balance().Int64(),
		Nonce:   int64(obj.Nonce()),
	}
}
