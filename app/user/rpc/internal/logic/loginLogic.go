package logic

import (
	"context"
	"fmt"

	"amigo-api/app/user/rpc/internal/logic/service/userService"
	"amigo-api/app/user/rpc/internal/svc"
	"amigo-api/common/pb"
	"amigo-api/common/utils"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *pb.LoginResp) (*pb.LoginSuccessResp, error) {
	user, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, in.Username)
	if err != nil {
		return nil, err
	}

	if user.UserId == 0 {
		return nil, fmt.Errorf("用户不存在")
	}

	pwd := utils.Md5(in.Password)
	if user.Password != pwd {
		return nil, fmt.Errorf("密码错误")
	}

	resp := &pb.LoginSuccessResp{}
	_ = copier.Copy(resp, user)
	resp.Token, err = userService.EncodeJwtToken(l.ctx, l.svcCtx, user.UserId)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
