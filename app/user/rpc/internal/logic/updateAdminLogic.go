package logic

import (
	"context"

	"amigo-api/app/user/model"

	"amigo-api/app/user/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAdminLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAdminLogic {
	return &UpdateAdminLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateAdminLogic) UpdateAdmin(in *pb.UpdateAdminReq) (*pb.SuccessResp, error) {
	// 检查数据是否存在
	_, err := l.svcCtx.AdminModel.FindOne(l.ctx, in.AdminId)
	if err != nil {
		if err == model.ErrNotFound {
			return &pb.SuccessResp{Success: false}, model.ErrNotFound
		}
		return &pb.SuccessResp{Success: false}, err
	}

	// 创建数据模型
	var m model.Admin
	if err := copier.Copy(&m, in); err != nil {
		return &pb.SuccessResp{Success: false}, err
	}

	// 检查重复
	isDuplicate, err := l.svcCtx.AdminModel.CheckDuplicate(l.ctx, &m)
	if err != nil {
		return &pb.SuccessResp{Success: false}, err
	}
	if isDuplicate {
		return &pb.SuccessResp{Success: false}, model.ErrDuplicate
	}

	// 更新数据
	if err := l.svcCtx.AdminModel.Update(l.ctx, &m); err != nil {
		return &pb.SuccessResp{Success: false}, err
	}

	return &pb.SuccessResp{Success: true}, nil
}
