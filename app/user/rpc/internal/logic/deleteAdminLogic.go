package logic

import (
	"context"

	"amigo-api/app/user/model"
	"amigo-api/app/user/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteAdminLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAdminLogic {
	return &DeleteAdminLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteAdminLogic) DeleteAdmin(in *pb.DeleteAdminReq) (*pb.SuccessResp, error) {
	// 检查主键id是否存在
	if in.AdminId == 0 {
		return &pb.SuccessResp{Success: false}, model.ErrNotFound
	}

	// 根据主键id删除数据
	if err := l.svcCtx.AdminModel.Delete(l.ctx, in.AdminId); err != nil {
		if err == model.ErrNotFound {
			return &pb.SuccessResp{Success: false}, model.ErrNotFound
		}
		return &pb.SuccessResp{Success: false}, err
	}

	return &pb.SuccessResp{Success: true}, nil
}
