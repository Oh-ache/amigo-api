package logic

import (
	"context"

	"amigo-api/app/baseCode/model"

	"amigo-api/app/baseCode/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateBaseCodeItemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateBaseCodeItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateBaseCodeItemLogic {
	return &UpdateBaseCodeItemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateBaseCodeItemLogic) UpdateBaseCodeItem(in *pb.BaseCodeItemResp) (*pb.BaseCodeItemResp, error) {
	// 检查数据是否存在
	_, err := l.svcCtx.BaseCodeItemModel.FindOne(l.ctx, in.BaseCodeItemId)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, model.ErrNotFound
		}
		return nil, err
	}

	// 创建数据模型
	var m model.BaseCodeItem
	if err := copier.Copy(&m, in); err != nil {
		l.Errorf("Failed to copy request data to model: %v", err)
		return nil, err
	}

	// 检查重复
	isDuplicate, err := l.svcCtx.BaseCodeItemModel.CheckDuplicate(l.ctx, &m)
	if err != nil {
		return nil, err
	}
	if isDuplicate {
		return nil, model.ErrDuplicate
	}

	// 更新数据
	if err := l.svcCtx.BaseCodeItemModel.Update(l.ctx, &m); err != nil {
		return nil, err
	}

	// 构造响应
	var resp pb.BaseCodeItemResp
	if err := copier.Copy(&resp, &m); err != nil {
		return nil, err
	}

	return &resp, nil
}
