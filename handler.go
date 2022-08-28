package main

import (
	"context"

	"github.com/MonteCarloClub/KBD/handler"
	"github.com/MonteCarloClub/KBD/kitex_gen/api"
)

// KanBanDatabaseImpl implements the last service interface defined in the IDL.
type KanBanDatabaseImpl struct{}

// GetData implements the KanBanDatabaseImpl interface.
func (s *KanBanDatabaseImpl) GetData(ctx context.Context, req *api.GetDataRequest) (resp *api.GetDataResponse, err error) {
	return handler.GetData(ctx, req)
}

// PutData implements the KanBanDatabaseImpl interface.
func (s *KanBanDatabaseImpl) PutData(ctx context.Context, req *api.PutDataRequest) (resp *api.PutDataResponse, err error) {
	return handler.PutData(ctx, req)
}

// GetAccountData implements the KanBanDatabaseImpl interface.
func (s *KanBanDatabaseImpl) GetAccountData(ctx context.Context, req *api.GetAccountDataRequest) (resp *api.GetAccountDataResponse, err error) {
	return handler.GetAccountData(ctx, req)
}
