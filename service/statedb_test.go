package service

import (
	"fmt"
	"path"
	"testing"

	"github.com/MonteCarloClub/KBD/model/kdb"
	"github.com/MonteCarloClub/KBD/model/state"

	"github.com/MonteCarloClub/KBD/common"
	"github.com/MonteCarloClub/KBD/constant"
)

func TestGetAccountData(t *testing.T) {
	file := path.Join("/", constant.DataDir, constant.StateDBFile)
	db, _ := kdb.NewLDBDatabase(file)
	statedb := state.New(common.Hash{}, db)
	obj := statedb.GetStateObject(common.StringToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d87"))
	fmt.Println(obj)
}
