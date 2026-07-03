package logic

import (
	"context"

	"amigo-api/app/device/model"
	"amigo-api/app/device/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteFirmwareLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteFirmwareLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFirmwareLogic {
	return &DeleteFirmwareLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteFirmwareLogic) DeleteFirmware(in *pb.DeleteFirmwareReq) (*pb.DeleteFirmwareResp, error) {
	_, err := l.svcCtx.FirmwareModel.FindOne(l.ctx, in.FirmwareId)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, model.ErrNotFound
		}
		return nil, err
	}

	if err := l.svcCtx.FirmwareModel.Delete(l.ctx, in.FirmwareId); err != nil {
		return nil, err
	}

	return &pb.DeleteFirmwareResp{Success: true}, nil
}
