package logic

import (
	"context"

	"amigo-api/app/baseCode/model"

	"amigo-api/app/baseCode/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateBaseCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateBaseCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateBaseCodeLogic {
	return &UpdateBaseCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateBaseCodeLogic) UpdateBaseCode(in *pb.BaseCodeResp) (*pb.BaseCodeResp, error) {
	// 检查数据是否存在
	_, err := l.svcCtx.BaseCodeModel.FindOne(l.ctx, in.BaseCodeId)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, model.ErrNotFound
		}
		return nil, err
	}

	// 创建数据模型
	var m model.BaseCode
	if err := copier.Copy(&m, in); err != nil {
		return nil, err
	}

	// 检查重复
	isDuplicate, err := l.svcCtx.BaseCodeModel.CheckDuplicate(l.ctx, &m)
	if err != nil {
		return nil, err
	}
	if isDuplicate {
		return nil, model.ErrDuplicate
	}

	// 更新数据
	if err := l.svcCtx.BaseCodeModel.Update(l.ctx, &m); err != nil {
		return nil, err
	}

	// 构造响应
	var resp pb.BaseCodeResp
	if err := copier.Copy(&resp, &m); err != nil {
		return nil, err
	}

	return &resp, nil
}
