package logic

import (
	"context"

	"amigo-api/app/user/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserThirdLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserThirdLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserThirdLoginLogic {
	return &UserThirdLoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserThirdLoginLogic) UserThirdLogin(in *pb.UserThirdLoginReq) (*pb.LoginSuccessResp, error) {
	// todo: add your logic here and delete this line

	return &pb.LoginSuccessResp{}, nil
}
