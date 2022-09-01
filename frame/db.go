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
	stateDB = state.New(common.Hash{31, 201, 87, 101, 23, 89, 84, 129, 135, 164, 77, 141, 1, 242, 108, 199, 216, 127, 167, 134, 81, 152, 56, 213, 89, 118, 4, 212, 140, 185, 223, 175}, GetDB())
}

func GetStateDB() *state.StateDB {
	if stateDB == nil {
		initStateDB()
	}
	return stateDB
}
