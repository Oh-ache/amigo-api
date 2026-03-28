package logic

import (
	"context"

	"amigo-api/app/sdk/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetWeatherLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetWeatherLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWeatherLogic {
	return &GetWeatherLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetWeatherLogic) GetWeather(in *pb.GetWeatherReq) (*pb.GetWeatherResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetWeatherResp{}, nil
}
