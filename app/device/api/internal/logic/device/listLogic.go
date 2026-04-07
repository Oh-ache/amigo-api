package device

import (
	"context"

	"amigo-api/app/device/api/internal/svc"
	"amigo-api/app/device/api/internal/types"
	"amigo-api/common/pb"
	"amigo-api/common/utils"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListLogic) List(req *types.ListDeviceReq) (resp *types.ListDeviceResp, err error) {
	resp = &types.ListDeviceResp{}
	param := &pb.ListDeviceReq{}

	copier.Copy(param, req)

	payload := &utils.JwtPayload{}
	_ = utils.DecodeJwtToken(l.ctx.Value("payload"), payload)

	if payload.Domain == "amigo-api" {
		param.UserId = payload.UserId
	}
	rpcResp, err := l.svcCtx.DeviceRpcClient.ListDevice(l.ctx, param)
	if err != nil {
		return nil, err
	}

	copier.Copy(resp, rpcResp)
	return resp, nil
}
