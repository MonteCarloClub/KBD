package frame

import (
	"path"

	"github.com/cloudwego/kitex/pkg/klog"

	"github.com/MonteCarloClub/KBD/common"
	"github.com/MonteCarloClub/KBD/constant"
	"github.com/MonteCarloClub/KBD/model/kdb"
	"github.com/MonteCarloClub/KBD/model/state"
)

var runState *state.StateDB
var stateDB *kdb.LDBDatabase
var blockDB *kdb.LDBDatabase
var root []byte

func Init() {
	initState()
	initBlock()
}

func initBlock() {
	file := path.Join("/", constant.DataDir, constant.BlockDBFile)
	blockDB, _ = kdb.NewLDBDatabase(file)
}

func GetRoot() []byte {
	if len(root) != 0 {
		return root
	}
	if blockDB == nil {
		initBlock()
	}
	res, err := blockDB.Get([]byte("root"))
	if err != nil {
		klog.Errorf("[GetRoot] get root failed %v", err)
	}
	return res
}
func PutRoot(value []byte) {
	if blockDB == nil {
		initBlock()
	}
	err := blockDB.Put([]byte("root"), value)
	if err != nil {
		klog.Info("[PutRoot] put root failed")
	}
	blockDB.Flush()
	return
}

func initStateDB() {
	file := path.Join("/", constant.DataDir, constant.StateDBFile)
	stateDB, _ = kdb.NewLDBDatabase(file)
}

func GetDB() *kdb.LDBDatabase {
	if stateDB == nil {
		initStateDB()
	}
	return stateDB
}

func initState() {
	runState = state.New(common.BytesToHash(GetRoot()), GetDB())
}

func GetState() *state.StateDB {
	if runState == nil {
		initState()
	}
	return runState
}
