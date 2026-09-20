package weather

import (
	"context"

	"amigo-api/app/sdk/api/internal/svc"
	"amigo-api/app/sdk/api/internal/types"
	"amigo-api/app/sdk/rpc/sdk"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type WeatherLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWeatherLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WeatherLogic {
	return &WeatherLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WeatherLogic) Weather(req *types.WeatherReq) (resp *types.WeatherResp, err error) {
	resp = &types.WeatherResp{}
	param := &sdk.GetWeatherReq{}

	copier.Copy(param, req)
	rpcResp, err := l.svcCtx.SdkRpcClient.GetWeather(l.ctx, param)
	if err != nil {
		return nil, err
	}

	copier.Copy(resp, rpcResp)
	return resp, nil
}