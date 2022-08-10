package main

import (
	"context"
	"github.com/MonteCarloClub/KBD/kitex_gen/api"
)

// KanBanDatabaseImpl implements the last service interface defined in the IDL.
type KanBanDatabaseImpl struct{}

// GetAccountData implements the KanBanDatabaseImpl interface.
func (s *KanBanDatabaseImpl) GetAccountData(ctx context.Context, req *api.GetAccountDataRequest) (resp *api.GetAccountDataResponse, err error) {
	// TODO: Your code here...
	return
}
