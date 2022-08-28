package handler

import (
	"context"
	"github.com/MonteCarloClub/KBD/kitex_gen/api"
	"github.com/MonteCarloClub/KBD/service"
)

// GetData implements the KanBanDatabaseImpl interface.
func GetData(ctx context.Context, req *api.GetDataRequest) (resp *api.GetDataResponse, err error) {
	resp = &api.GetDataResponse{}
	value, err := service.GetData(ctx, req)
	if err != nil {
		return resp, err
	}
	resp.Value = string(value)
	return resp, nil
}

// PutData implements the KanBanDatabaseImpl interface.
func PutData(ctx context.Context, req *api.PutDataRequest) (resp *api.PutDataResponse, err error) {
	resp = &api.PutDataResponse{}
	err = service.PutData(ctx, req)
	if err == nil {
		resp.Success = true
	}
	return resp, err
}
