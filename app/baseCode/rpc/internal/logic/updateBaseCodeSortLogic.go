package logic

import (
	"context"

	"amigo-api/app/baseCode/model"

	"amigo-api/app/baseCode/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateBaseCodeSortLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateBaseCodeSortLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateBaseCodeSortLogic {
	return &UpdateBaseCodeSortLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateBaseCodeSortLogic) UpdateBaseCodeSort(in *pb.BaseCodeSortResp) (*pb.BaseCodeSortResp, error) {
	// 创建 model.BaseCodeSort 实例
	var sort model.BaseCodeSort

	// 使用 copier 库将入参复制到 sort 实例
	if err := copier.Copy(&sort, in); err != nil {
		return nil, err
	}

	// 先判断是否重复
	isDuplicate, err := l.svcCtx.BaseCodeSortModel.CheckDuplicate(l.ctx, &sort)
	if err != nil {
		return nil, err
	}
	if isDuplicate {
		return nil, model.ErrDuplicate
	}

	// 更新数据
	if err := l.svcCtx.BaseCodeSortModel.Update(l.ctx, &sort); err != nil {
		return nil, err
	}

	// 查询更新后的完整数据
	updatedSort, err := l.svcCtx.BaseCodeSortModel.FindOne(l.ctx, sort.BaseCodeSortId)
	if err != nil {
		return nil, err
	}

	// 创建返回响应
	var resp pb.BaseCodeSortResp
	if err := copier.Copy(&resp, updatedSort); err != nil {
		return nil, err
	}

	return &resp, nil
}
