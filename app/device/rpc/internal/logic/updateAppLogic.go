package logic

import (
	"context"

	"amigo-api/app/device/model"
	"amigo-api/app/device/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAppLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAppLogic {
	return &UpdateAppLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateAppLogic) UpdateApp(in *pb.AppResp) (*pb.AppResp, error) {
	_, err := l.svcCtx.AppModel.FindOne(l.ctx, in.AppId)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, model.ErrNotFound
		}
		return nil, err
	}
	var m model.App
	if err := copier.Copy(&m, in); err != nil {
		return nil, err
	}
	isDuplicate, err := l.svcCtx.AppModel.CheckDuplicate(l.ctx, &m)
	if err != nil {
		return nil, err
	}
	if isDuplicate {
		return nil, model.ErrDuplicate
	}
	if err := l.svcCtx.AppModel.Update(l.ctx, &m); err != nil {
		return nil, err
	}
	var resp pb.AppResp
	if err := copier.Copy(&resp, &m); err != nil {
		return nil, err
	}
	return &resp, nil
}
