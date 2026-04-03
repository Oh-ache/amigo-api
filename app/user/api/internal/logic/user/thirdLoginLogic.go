package user

import (
	"context"
	"fmt"

	"amigo-api/app/user/api/internal/svc"
	"amigo-api/app/user/api/internal/types"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type ThirdLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewThirdLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ThirdLoginLogic {
	return &ThirdLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ThirdLoginLogic) ThirdLogin(req *types.UserThirdLoginReq) (resp *types.UserLoginResp, err error) {
	resp = &types.UserLoginResp{}

	param := &pb.UserThirdLoginReq{
		AppType: req.AppType,
		Code:    req.Code,
	}
	fmt.Println("==============")
	fmt.Println(param)

	rpcResp, err := l.svcCtx.UserRpcClient.UserThirdLogin(l.ctx, param)
	if err != nil {
		return nil, err
	}

	copier.Copy(resp, rpcResp)
	return resp, nil
}
