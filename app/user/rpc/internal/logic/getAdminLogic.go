package logic

import (
	"context"

	"amigo-api/app/user/model"

	"amigo-api/app/user/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetAdminLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAdminLogic {
	return &GetAdminLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAdminLogic) GetAdmin(in *pb.GetAdminReq) (*pb.AdminResp, error) {
	// 根据主键id查询数据
	if in.AdminId != 0 {
		data, err := l.svcCtx.AdminModel.FindOne(l.ctx, in.AdminId)
		if err != nil {
			if err == model.ErrNotFound {
				return nil, model.ErrNotFound
			}
			return nil, err
		}
		var resp pb.AdminResp
		if err := copier.Copy(&resp, data); err != nil {
			return nil, err
		}
		return &resp, nil
	}

	// 根据username查询数据
	if in.Username != "" {
		data, err := l.svcCtx.AdminModel.FindOneByUsername(l.ctx, in.Username)
		if err != nil {
			if err == model.ErrNotFound {
				return nil, model.ErrNotFound
			}
			return nil, err
		}
		var resp pb.AdminResp
		if err := copier.Copy(&resp, data); err != nil {
			return nil, err
		}
		return &resp, nil
	}

	// 根据mobile查询数据
	if in.Mobile != "" {
		data, err := l.svcCtx.AdminModel.FindOneByMobile(l.ctx, in.Mobile)
		if err != nil {
			if err == model.ErrNotFound {
				return nil, model.ErrNotFound
			}
			return nil, err
		}
		var resp pb.AdminResp
		if err := copier.Copy(&resp, data); err != nil {
			return nil, err
		}
		return &resp, nil
	}

	return nil, model.ErrNotFound
}
