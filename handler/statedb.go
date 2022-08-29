package handler

import (
	"context"
	"fmt"

	"github.com/MonteCarloClub/KBD/kitex_gen/api"
	"github.com/MonteCarloClub/KBD/model"
	"github.com/MonteCarloClub/KBD/service"
	"github.com/MonteCarloClub/KBD/util"
	"github.com/cloudwego/kitex/pkg/klog"
)

// GetAccountData implements the KanBanDatabaseImpl interface.
func GetAccountData(ctx context.Context, req *api.GetAccountDataRequest) (resp *api.GetAccountDataResponse, err error) {
	resp = &api.GetAccountDataResponse{}
	if req.Address == "" {
		return nil, fmt.Errorf("wrong account")
	}
	obj := service.GetAccountData(ctx, req.GetAddress())
	resp.Account = model.StateObject2VO(obj)
	klog.CtxInfof(ctx, "[GetAccountData]req = %v,resp = %v", util.ToString(req), util.ToString(resp))
	return resp, nil
}
