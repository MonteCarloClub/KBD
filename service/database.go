package service

import (
	"context"

	"github.com/MonteCarloClub/KBD/frame"
	"github.com/MonteCarloClub/KBD/kitex_gen/api"
)

// GetData implements the KanBanDatabaseImpl interface.
func GetData(ctx context.Context, req *api.GetDataRequest) (valuse []byte, err error) {
	return frame.GetDB().Get([]byte(req.Key))
}

// PutData implements the KanBanDatabaseImpl interface.
func PutData(ctx context.Context, req *api.PutDataRequest) (err error) {
	err = frame.GetDB().Put([]byte(req.Key), []byte(req.Value))
	return err
}
