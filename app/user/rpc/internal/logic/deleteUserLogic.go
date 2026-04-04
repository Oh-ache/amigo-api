package logic

import (
	"context"

	"amigo-api/app/user/model"
	"amigo-api/app/user/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserLogic {
	return &DeleteUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteUserLogic) DeleteUser(in *pb.DeleteUserReq) (*pb.SuccessResp, error) {
	// 检查主键id是否存在
	if in.UserId == 0 {
		return &pb.SuccessResp{Success: false}, model.ErrNotFound
	}

	// 根据主键id删除数据
	if err := l.svcCtx.UserModel.Delete(l.ctx, in.UserId); err != nil {
		if err == model.ErrNotFound {
			return &pb.SuccessResp{Success: false}, model.ErrNotFound
		}
		return &pb.SuccessResp{Success: false}, err
	}

	return &pb.SuccessResp{Success: true}, nil
}
