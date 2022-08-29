package frame

import (
	"context"
	"path"

	"github.com/MonteCarloClub/KBD/model/kdb"
	"github.com/MonteCarloClub/KBD/model/state"

	"github.com/MonteCarloClub/KBD/common"
	"github.com/MonteCarloClub/KBD/constant"
)

func Init(ctx context.Context) {
	initStateDB(ctx)
}

func initStateDB(ctx context.Context) {
	file := path.Join("/", constant.DataDir, constant.StateDBFile)
	db, _ := kdb.NewLDBDatabase(file)
	statedb := state.New(common.Hash{}, db)
	context.WithValue(ctx, constant.CtxStateDB, statedb)
}

func GetStateDB(ctx context.Context) *state.StateDB {
	value := ctx.Value(constant.CtxStateDB)
	if value != nil {
		return value.(*state.StateDB)
	}
	return nil
}
