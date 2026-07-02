package logic

import (
	"context"

	"amigo-api/app/device/model"
	"amigo-api/app/device/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddAppLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddAppLogic {
	return &AddAppLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddAppLogic) AddApp(in *pb.AddAppReq) (*pb.AppResp, error) {
	var m model.App
	if err := copier.Copy(&m, in); err != nil {
		return nil, err
	}
	if m.IsDelete == 0 {
		m.IsDelete = 2
	}
	isDuplicate, err := l.svcCtx.AppModel.CheckDuplicate(l.ctx, &m)
	if err != nil {
		return nil, err
	}
	if isDuplicate {
		existing, err := l.svcCtx.AppModel.FindOneByAppKey(l.ctx, m.AppKey)
		if err != nil {
			return nil, err
		}
		m.AppId = existing.AppId
		if err := l.svcCtx.AppModel.Update(l.ctx, &m); err != nil {
			return nil, err
		}
	} else {
		result, err := l.svcCtx.AppModel.Insert(l.ctx, &m)
		if err != nil {
			return nil, err
		}
		id, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}
		m.AppId = uint64(id)
	}
	var resp pb.AppResp
	if err := copier.Copy(&resp, &m); err != nil {
		return nil, err
	}
	return &resp, nil
}
