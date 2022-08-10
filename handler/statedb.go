package handler

import (
	"context"
	"fmt"
	"github.com/MonteCarloClub/KBD/kitex_gen/api"
	"github.com/MonteCarloClub/KBD/model"
	"github.com/MonteCarloClub/KBD/service"
)

// GetAccountData implements the KanBanDatabaseImpl interface.
func GetAccountData(ctx context.Context, req *api.GetAccountDataRequest) (resp *api.GetAccountDataResponse, err error) {
	if req.Address == "" {
		return nil, fmt.Errorf("wrong account")
	}
	obj := service.GetAccountData(ctx, req.GetAddress())
	resp.Account = model.StateObject2VO(obj)
	return resp, nil
}
