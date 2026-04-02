package logic

import (
	"context"
	"fmt"

	"amigo-api/app/user/rpc/internal/logic/service/adminService"
	"amigo-api/app/user/rpc/internal/svc"
	"amigo-api/common/pb"
	"amigo-api/common/utils"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginAdminLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginAdminLogic {
	return &LoginAdminLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginAdminLogic) LoginAdmin(in *pb.LoginAdminReq) (*pb.AdminLoginSuccessResp, error) {
	admin, err := l.svcCtx.AdminModel.FindOneByMobile(l.ctx, in.Mobile)
	if err != nil {
		return nil, err
	}

	if admin.AdminId == 0 {
		return nil, fmt.Errorf("用户不存在")
	}

	pwd := utils.Md5(in.Password)
	if admin.Password != pwd {
		return nil, fmt.Errorf("密码错误")
	}

	resp := &pb.AdminLoginSuccessResp{}
	_ = copier.Copy(resp, admin)
	resp.Token, err = adminService.EncodeJwtToken(l.ctx, l.svcCtx, admin.AdminId)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
