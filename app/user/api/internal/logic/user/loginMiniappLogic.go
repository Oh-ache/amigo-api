package user

import (
	"context"

	"amigo-api/app/user/api/internal/svc"
	"amigo-api/app/user/api/internal/types"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginMiniappLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginMiniappLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginMiniappLogic {
	return &LoginMiniappLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginMiniappLogic) LoginMiniapp(req *types.UserLoginMiniappReq) (resp *types.UserLoginResp, err error) {
	resp = &types.UserLoginResp{}

	param := &pb.UserThirdLoginReq{
		AppType: req.AppType,
		Code:    req.Code,
	}

	rpcResp, err := l.svcCtx.UserRpcClient.UserThirdLogin(l.ctx, param)
	if err != nil {
		return nil, err
	}

	copier.Copy(resp, rpcResp)
	return resp, nil
}
