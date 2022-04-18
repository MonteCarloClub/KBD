package kdb

import (
	"github.com/MonteCarloClub/KBD/common"
	"os"
	"path"
	"testing"
)

func TestNewDb(t *testing.T) {
	file := path.Join("/", "tmp", "ldbtesttmpfile")
	if common.FileExist(file) {
		os.RemoveAll(file)
	}
	db, _ := NewLDBDatabase(file)
	db.Put([]byte("dfdsaf"), []byte("dddd"))
}
