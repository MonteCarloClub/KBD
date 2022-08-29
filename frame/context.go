package frame

import (
	"path"

	"github.com/MonteCarloClub/KBD/common"
	"github.com/MonteCarloClub/KBD/constant"
	"github.com/MonteCarloClub/KBD/model/kdb"
	"github.com/MonteCarloClub/KBD/model/state"
)

var stateDB *state.StateDB
var db *kdb.LDBDatabase

func Init() {
	initDB()
	initStateDB()
}

func initDB() {
	file := path.Join("/", constant.DataDir, constant.StateDBFile)
	db, _ = kdb.NewLDBDatabase(file)
}

func GetDB() *kdb.LDBDatabase {
	if db == nil {
		initDB()
	}
	return db
}

func initStateDB() {
	stateDB = state.New(common.Hash{}, GetDB())
}

func GetStateDB() *state.StateDB {
	if stateDB == nil {
		initStateDB()
	}
	return stateDB
}
