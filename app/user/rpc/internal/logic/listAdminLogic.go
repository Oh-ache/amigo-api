package logic

import (
	"context"

	"amigo-api/app/user/model"
	"amigo-api/app/user/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListAdminLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAdminLogic {
	return &ListAdminLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListAdminLogic) ListAdmin(in *pb.ListAdminReq) (*pb.ListAdminResp, error) {
	// 构建查询条件
	search := &model.AdminSearch{}
	if err := copier.Copy(search, in); err != nil {
		return nil, err
	}

	// 查询数据
	list, total, err := l.svcCtx.AdminModel.List(l.ctx, search)
	if err != nil {
		return nil, err
	}

	// 构造响应
	resp := &pb.ListAdminResp{}
	if err := copier.Copy(resp, &struct {
		List  []*model.Admin
		Total int64
	}{
		List:  list,
		Total: total,
	}); err != nil {
		return nil, err
	}

	return resp, nil
}
