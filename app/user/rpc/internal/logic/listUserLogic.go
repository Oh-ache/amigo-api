package logic

import (
	"context"

	"amigo-api/app/user/model"
	"amigo-api/app/user/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUserLogic {
	return &ListUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListUserLogic) ListUser(in *pb.ListUserReq) (*pb.ListUserResp, error) {
	// 构建查询条件
	search := &model.UserSearch{}
	if err := copier.Copy(search, in); err != nil {
		return nil, err
	}

	// 查询数据
	list, total, err := l.svcCtx.UserModel.List(l.ctx, search)
	if err != nil {
		return nil, err
	}

	// 构造响应
	resp := &pb.ListUserResp{}
	if err := copier.Copy(resp, &struct {
		List  []*model.User
		Total int64
	}{
		List:  list,
		Total: total,
	}); err != nil {
		return nil, err
	}

	return resp, nil
}
