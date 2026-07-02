package logic

import (
	"context"
	"fmt"

	"amigo-api/app/user/rpc/internal/logic/service/adminService"
	"amigo-api/app/user/rpc/internal/svc"
	"amigo-api/common/pb"
	"amigo-api/common/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshAdminTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRefreshAdminTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshAdminTokenLogic {
	return &RefreshAdminTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RefreshAdminTokenLogic) RefreshAdminToken(in *pb.RefreshTokenReq) (*pb.RefreshTokenResp, error) {
	if in.RefreshToken == "" {
		return nil, fmt.Errorf("refresh_token 不能为空")
	}

	// 验签 + 解析 payload
	var payload utils.JwtPayload
	if err := utils.DecodeJwtTokenFromString(l.svcCtx.Config.JwtAuth.AccessSecret, in.RefreshToken, &payload); err != nil {
		return nil, fmt.Errorf("无效的 refresh_token: %v", err)
	}

	// 必须为 refresh 类型（防 access token 误用）
	if payload.TokenType != utils.TokenTypeRefresh {
		return nil, fmt.Errorf("token 类型错误，期望 refresh")
	}

	if payload.UserId == 0 {
		return nil, fmt.Errorf("user_id 缺失")
	}

	// 校验用户存在
	admin, err := l.svcCtx.AdminModel.FindOne(l.ctx, payload.UserId)
	if err != nil {
		return nil, fmt.Errorf("用户不存在")
	}
	if admin.AdminId == 0 {
		return nil, fmt.Errorf("用户不存在")
	}

	resp := &pb.RefreshTokenResp{}
	resp.AccessToken, resp.ExpiresIn, err = adminService.EncodeAccessToken(l.ctx, l.svcCtx, admin.AdminId)
	if err != nil {
		return nil, err
	}
	resp.RefreshToken, _, err = adminService.EncodeRefreshToken(l.ctx, l.svcCtx, admin.AdminId)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
