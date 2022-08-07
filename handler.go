package main

import (
	"os"
	"path"

	"github.com/MonteCarloClub/KBD/common"
	"github.com/MonteCarloClub/KBD/kdb"
	"github.com/MonteCarloClub/KBD/rpcserver"
	"github.com/MonteCarloClub/KBD/state"
	"github.com/tidwall/gjson"
)

// handleGetBlockHash implements the getblockhash command.
func GetFunc(s *rpcserver.RpcServer, cmd interface{}, closeChan <-chan struct{}) (interface{}, error) {
	res := gjson.Parse(cmd.(string))
	if res.Get("funcName").Value().(string) == "GetAccountData" {
		return GetAccountData(res.Get("value"))
	}
	return nil, nil
}

func GetAccountData(res gjson.Result) (interface{}, error) {
	file := path.Join("/", "tmp", "kdb")
	if common.FileExist(file) {
		os.RemoveAll(file)
	}
	db, _ := kdb.NewLDBDatabase(file)
	addr := res.Value().(common.Address)
	state := state.New(common.Hash{}, db)
	return state.GetStateObject(addr), nil
}

func GetBlockData(s *rpcserver.RpcServer, cmd interface{}, closeChan <-chan struct{}) (interface{}, error) {
	// 这一行测试用
	return cmd, nil
}
