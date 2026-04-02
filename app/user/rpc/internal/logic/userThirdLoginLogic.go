package logic

import (
	"context"
	"fmt"

	"amigo-api/app/user/rpc/internal/logic/service/userService"
	"amigo-api/app/user/rpc/internal/logic/service/userThirdService"
	"amigo-api/app/user/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
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
	session := &pb.MiniappCodeResp{}
	var err error
	if session, err = userThirdService.Code2Session(l.ctx, l.svcCtx, in.AppType, in.Code); err != nil {
		return nil, err
	}

	if session.Errcode != 0 {
		return nil, fmt.Errorf("认证失败")
	}

	// 自动注册获取用户信息
	user, err := userThirdService.InsertUser(l.ctx, l.svcCtx, in.AppType, session.Openid, session.Unionid)
	if err != nil {
		return nil, err
	}

	resp := &pb.LoginSuccessResp{}
	_ = copier.Copy(resp, user)
	resp.Token, err = userService.EncodeJwtToken(l.ctx, l.svcCtx, user.UserId)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
