package logic

import (
	"context"

	"amigo-api/app/user/model"

	"amigo-api/app/user/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserLogic) UpdateUser(in *pb.UpdateUserReq) (*pb.SuccessResp, error) {
	// 检查数据是否存在
	_, err := l.svcCtx.UserModel.FindOne(l.ctx, in.UserId)
	if err != nil {
		if err == model.ErrNotFound {
			return &pb.SuccessResp{Success: false}, model.ErrNotFound
		}
		l.Errorf("Failed to find user by id %d: %v", in.UserId, err)
		return &pb.SuccessResp{Success: false}, err
	}

	// 创建数据模型
	var m model.User
	if err := copier.Copy(&m, in); err != nil {
		return &pb.SuccessResp{Success: false}, err
	}

	// 检查重复
	isDuplicate, err := l.svcCtx.UserModel.CheckDuplicate(l.ctx, &m)
	if err != nil {
		return &pb.SuccessResp{Success: false}, err
	}
	if isDuplicate {
		return &pb.SuccessResp{Success: false}, model.ErrDuplicate
	}

	// 更新数据
	if err := l.svcCtx.UserModel.Update(l.ctx, &m); err != nil {
		return &pb.SuccessResp{Success: false}, err
	}

	return &pb.SuccessResp{Success: true}, nil
}
