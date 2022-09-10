package service

import (
	"fmt"
	"testing"

	"github.com/MonteCarloClub/KBD/common"
	"github.com/MonteCarloClub/KBD/frame"
)

func TestPutRoot(t *testing.T) {
	frame.PutRoot([]byte{31, 201, 87, 101, 23, 89, 84, 129, 135, 164, 77, 141, 1, 242, 108, 199, 216, 127, 167, 134, 81, 152, 56, 213, 89, 118, 4, 212, 140, 185, 223, 175})
}

func TestPutAccountData(t *testing.T) {
	stateDB := frame.GetState()
	address := common.HexToAddress("0x945304eb96065b2a98b57a48a06ae28d285a71b4")
	balance := "0x17"
	code := "0x6000355415600957005b60203560003555"
	nonce := "0x00"
	obj := StateObjectFromAccount(frame.GetDB(), address, balance, code, nonce)
	stateDB.UpdateStateObject(obj)
	stateDB.Sync()
	frame.PutRoot(stateDB.Trie().Root())
	res := stateDB.GetStateObject(address)
	fmt.Println(res)
}

func TestGetAccountData(t *testing.T) {
	stateDB := frame.GetState()
	address := common.HexToAddress("0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b")
	res := stateDB.GetStateObject(address)
	fmt.Println(res)
}
