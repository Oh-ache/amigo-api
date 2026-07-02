package logic

import (
	"context"

	"amigo-api/app/device/model"
	"amigo-api/app/device/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetAppLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAppLogic {
	return &GetAppLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAppLogic) GetApp(in *pb.GetAppReq) (*pb.AppResp, error) {
	var app *model.App
	var err error
	if in.AppId != 0 {
		app, err = l.svcCtx.AppModel.FindOne(l.ctx, in.AppId)
	} else if in.AppKey != "" {
		app, err = l.svcCtx.AppModel.FindOneByAppKey(l.ctx, in.AppKey)
	} else {
		return nil, model.ErrNotFound
	}
	if err != nil {
		if err == model.ErrNotFound {
			return nil, model.ErrNotFound
		}
		return nil, err
	}
	var resp pb.AppResp
	if err := copier.Copy(&resp, app); err != nil {
		return nil, err
	}
	return &resp, nil
}
