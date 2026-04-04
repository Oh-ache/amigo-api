package logic

import (
	"context"

	"amigo-api/app/baseCode/model"
	"amigo-api/app/baseCode/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListBaseCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListBaseCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListBaseCodeLogic {
	return &ListBaseCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListBaseCodeLogic) ListBaseCode(in *pb.ListBaseCodeReq) (*pb.ListBaseCodeResp, error) {
	// 构建查询条件
	search := &model.BaseCodeSearch{}
	if err := copier.Copy(search, in); err != nil {
		return nil, err
	}

	// 查询数据
	list, total, err := l.svcCtx.BaseCodeModel.List(l.ctx, search)
	if err != nil {
		return nil, err
	}

	// 构造响应
	resp := &pb.ListBaseCodeResp{}
	if err := copier.Copy(resp, &struct {
		List  []*model.BaseCode
		Total int64
	}{
		List:  list,
		Total: total,
	}); err != nil {
		return nil, err
	}

	return resp, nil
}
