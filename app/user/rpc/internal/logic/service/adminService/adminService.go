package adminService

import (
	"context"
	"fmt"
	"time"

	"amigo-api/app/user/model"
	"amigo-api/app/user/rpc/internal/svc"
	"amigo-api/common/utils"
)

func CheckDuplicate(ctx context.Context, svcCtx *svc.ServiceContext, admin *model.Admin) error {
	exists, err := svcCtx.AdminModel.CheckDuplicate(ctx, admin)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("管理员已存在")
	}

	return nil
}

func Insert(ctx context.Context, svcCtx *svc.ServiceContext, admin *model.Admin, nodata bool) (*model.Admin, error) {
	// 插入管理员数据
	_, err := svcCtx.AdminModel.Insert(ctx, admin)
	if err != nil {
		return nil, err
	}

	// 不返回管理员数据
	if nodata {
		return nil, nil
	}

	// 查询数据
	resp, err := svcCtx.AdminModel.FindOneByUsername(ctx, admin.Username)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// EncodeAccessToken 生成短期 access token（默认 2 小时）
func EncodeAccessToken(ctx context.Context, svcCtx *svc.ServiceContext, userId uint64) (token string, expiresAt int64, err error) {
	return encodeToken(svcCtx, userId, utils.TokenTypeAccess, svcCtx.Config.JwtAuth.AccessExpire)
}

// EncodeRefreshToken 生成长期 refresh token（1 天）
func EncodeRefreshToken(ctx context.Context, svcCtx *svc.ServiceContext, userId uint64) (token string, expiresAt int64, err error) {
	return encodeToken(svcCtx, userId, utils.TokenTypeRefresh, utils.RefreshExpire)
}

func encodeToken(svcCtx *svc.ServiceContext, userId uint64, tokenType string, expire int64) (string, int64, error) {
	payload := &utils.JwtPayload{
		UserId:    userId,
		Domain:    "amigo-admin",
		TokenType: tokenType,
	}
	now := time.Now().Unix()
	token, err := utils.EncodeJwtToken(
		svcCtx.Config.JwtAuth.AccessSecret,
		now,
		expire,
		payload,
	)
	if err != nil {
		return "", 0, fmt.Errorf("获取token失败")
	}
	return token, now + expire, nil
}

// EncodeJwtToken 兼容旧调用（默认 access token）
func EncodeJwtToken(ctx context.Context, svcCtx *svc.ServiceContext, userId uint64) (token string, err error) {
	token, _, err = EncodeAccessToken(ctx, svcCtx, userId)
	return
}
