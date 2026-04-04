package logic

import (
	"context"

	"amigo-api/app/baseCode/model"

	"amigo-api/app/baseCode/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListBaseCodeItemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListBaseCodeItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListBaseCodeItemLogic {
	return &ListBaseCodeItemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListBaseCodeItemLogic) ListBaseCodeItem(in *pb.ListBaseCodeItemReq) (*pb.ListBaseCodeItemResp, error) {
	// 构造查询条件
	search := &model.BaseCodeItemSearch{}
	_ = copier.Copy(search, in)

	// 调用模型方法查询数据
	list, total, err := l.svcCtx.BaseCodeItemModel.List(l.ctx, search)
	if err != nil {
		return nil, err
	}

	// 转换模型数据到响应格式
	var resp pb.ListBaseCodeItemResp
	resp.Total = total
	resp.List = make([]*pb.BaseCodeItemResp, 0, len(list))
	for _, item := range list {
		var itemResp pb.BaseCodeItemResp
		if err := copier.Copy(&itemResp, item); err != nil {
			continue
		}
		resp.List = append(resp.List, &itemResp)
	}

	return &resp, nil
}
