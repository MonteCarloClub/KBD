package service

import (
	"fmt"
	"testing"

	"github.com/MonteCarloClub/KBD/common"
	"github.com/MonteCarloClub/KBD/frame"
)

func TestPutAccountData(t *testing.T) {
	stateDB := frame.GetStateDB()
	address := common.StringToAddress("0x945304eb96065b2a98b57a48a06ae28d285a71b5")
	balance := "0x17"
	code := "0x6000355415600957005b60203560003555"
	nonce := "0x00"
	obj := StateObjectFromAccount(frame.GetDB(), address, balance, code, nonce)
	stateDB.UpdateStateObject(obj)
	stateDB.Trie().Commit()
	res := stateDB.GetStateObject(address)
	fmt.Println(res)
}

func TestGetAccountData(t *testing.T) {
	stateDB := frame.GetStateDB()
	address := common.StringToAddress("0x945304eb96065b2a98b57a48a06ae28d285a71b5")
	res := stateDB.GetStateObject(address)
	fmt.Println(res)
}
