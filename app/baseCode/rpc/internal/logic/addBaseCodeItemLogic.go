package logic

import (
	"context"

	"amigo-api/app/baseCode/model"

	"amigo-api/app/baseCode/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddBaseCodeItemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddBaseCodeItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddBaseCodeItemLogic {
	return &AddBaseCodeItemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddBaseCodeItemLogic) AddBaseCodeItem(in *pb.AddBaseCodeItemReq) (*pb.BaseCodeItemResp, error) {
	var m model.BaseCodeItem
	if err := copier.Copy(&m, in); err != nil {
		l.Errorf("Failed to copy request data to model: %v", err)
		return nil, err
	}

	if _, err := l.svcCtx.BaseCodeItemModel.Insert(l.ctx, &m); err != nil {
		l.Errorf("Failed to insert base code item: %v", err)
		return nil, err
	}

	var resp pb.BaseCodeItemResp
	if err := copier.Copy(&resp, &m); err != nil {
		l.Errorf("Failed to copy model to response: %v", err)
		return nil, err
	}

	return &resp, nil
}
