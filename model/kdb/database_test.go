package kdb

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/MonteCarloClub/KBD/common"
	"github.com/MonteCarloClub/KBD/constant"
)

func TestNewDb(t *testing.T) {
	file := path.Join("/", constant.DataDir, constant.StateDBFile)
	if common.FileExist(file) {
		os.RemoveAll(file)
	}
	db, _ := NewLDBDatabase(file)
	db.Put([]byte("dfdsaf"), []byte("asdjkfghkl;awes"))
	res, err := db.Get([]byte("dfdsaf"))
	fmt.Println(string(res), err)
	db.Close()
	db, _ = NewLDBDatabase(file)
	res, err = db.Get([]byte("dfdsaf"))
	fmt.Println(string(res), err)
}
