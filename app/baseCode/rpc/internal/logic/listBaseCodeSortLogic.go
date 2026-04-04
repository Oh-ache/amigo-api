package logic

import (
	"context"

	"amigo-api/app/baseCode/model"

	"amigo-api/app/baseCode/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListBaseCodeSortLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListBaseCodeSortLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListBaseCodeSortLogic {
	return &ListBaseCodeSortLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListBaseCodeSortLogic) ListBaseCodeSort(in *pb.ListBaseCodeSortReq) (*pb.ListBaseCodeSortResp, error) {
	// 构造查询条件
	search := &model.BaseCodeSortSearch{}
	_ = copier.Copy(search, in)

	// 调用模型方法查询数据
	list, total, err := l.svcCtx.BaseCodeSortModel.List(l.ctx, search)
	if err != nil {
		return nil, err
	}

	// 转换模型数据到响应格式
	var resp pb.ListBaseCodeSortResp
	resp.Total = total
	resp.List = make([]*pb.BaseCodeSortResp, 0, len(list))
	for _, sort := range list {
		var sortResp pb.BaseCodeSortResp
		if err := copier.Copy(&sortResp, sort); err != nil {
			continue
		}
		resp.List = append(resp.List, &sortResp)
	}

	return &resp, nil
}
