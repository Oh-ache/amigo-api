package logic

import (
	"context"

	"amigo-api/app/device/model"
	"amigo-api/app/device/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteAppLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAppLogic {
	return &DeleteAppLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteAppLogic) DeleteApp(in *pb.DeleteAppReq) (*pb.DeleteAppResp, error) {
	var app *model.App
	var err error
	if in.AppId != 0 {
		app, err = l.svcCtx.AppModel.FindOne(l.ctx, in.AppId)
	} else if in.AppKey != "" {
		app, err = l.svcCtx.AppModel.FindOneByAppKey(l.ctx, in.AppKey)
	} else {
		return &pb.DeleteAppResp{Success: false}, nil
	}
	if err != nil {
		if err == model.ErrNotFound {
			return &pb.DeleteAppResp{Success: false}, nil
		}
		return nil, err
	}
	if err := l.svcCtx.AppModel.Delete(l.ctx, app.AppId); err != nil {
		return nil, err
	}
	return &pb.DeleteAppResp{Success: true}, nil
}
