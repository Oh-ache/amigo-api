package logic

import (
	"context"

	"amigo-api/app/user/model"

	"amigo-api/app/user/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserLogic) GetUser(in *pb.GetUserReq) (*pb.UserResp, error) {
	// 根据主键id查询数据
	if in.UserId != 0 {
		data, err := l.svcCtx.UserModel.FindOne(l.ctx, in.UserId)
		if err != nil {
			if err == model.ErrNotFound {
				return nil, model.ErrNotFound
			}
			return nil, err
		}
		var resp pb.UserResp
		if err := copier.Copy(&resp, data); err != nil {
			return nil, err
		}
		return &resp, nil
	}

	// 根据username查询数据
	if in.Username != "" {
		data, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, in.Username)
		if err != nil {
			if err == model.ErrNotFound {
				return nil, model.ErrNotFound
			}
			return nil, err
		}
		var resp pb.UserResp
		if err := copier.Copy(&resp, data); err != nil {
			return nil, err
		}
		return &resp, nil
	}

	return nil, model.ErrNotFound
}
