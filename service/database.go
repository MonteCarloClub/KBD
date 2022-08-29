package service

import (
	"context"
	"path"

	"github.com/MonteCarloClub/KBD/constant"
	"github.com/MonteCarloClub/KBD/kdb"
	"github.com/MonteCarloClub/KBD/kitex_gen/api"
)

// GetData implements the KanBanDatabaseImpl interface.
func GetData(ctx context.Context, req *api.GetDataRequest) (valuse []byte, err error) {
	file := path.Join("/", constant.DataDir, constant.StateDBFile)
	db, _ := kdb.NewLDBDatabase(file)
	defer db.Close()
	return db.Get([]byte(req.Key))
}

// PutData implements the KanBanDatabaseImpl interface.
func PutData(ctx context.Context, req *api.PutDataRequest) (err error) {
	file := path.Join("/", constant.DataDir, constant.StateDBFile)
	db, _ := kdb.NewLDBDatabase(file)
	defer db.Close()
	err = db.Put([]byte(req.Key), []byte(req.Value))
	return err
}
