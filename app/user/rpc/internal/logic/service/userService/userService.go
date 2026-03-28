package userService

import (
	"context"
	"fmt"
	"time"

	"amigo-api/app/user/model"
	"amigo-api/app/user/rpc/internal/svc"
	"amigo-api/common/utils"
)

func EncodeJwtToken(ctx context.Context, svcCtx *svc.ServiceContext, userId uint64) (token string, err error) {
	jwtPayload := &utils.JwtPayload{
		UserId: userId,
		Domain: "amigo-api",
	}
	if token, err = utils.EncodeJwtToken(
		svcCtx.Config.JwtAuth.AccessSecret,
		time.Now().Unix(),
		svcCtx.Config.JwtAuth.AccessExpire,
		jwtPayload,
	); err != nil {
		return "", fmt.Errorf("获取token失败")
	}

	return
}

func CheckDuplicate(ctx context.Context, svcCtx *svc.ServiceContext, user *model.User) error {
	exists, err := svcCtx.UserModel.CheckDuplicate(ctx, user)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("用户已存在")
	}

	return nil
}

func Insert(ctx context.Context, svcCtx *svc.ServiceContext, user *model.User, nodata bool) (*model.User, error) {
	// 插入管理员数据
	_, err := svcCtx.UserModel.Insert(ctx, user)
	if err != nil {
		return nil, err
	}

	// 不返回管理员数据
	if nodata {
		return nil, nil
	}

	// 查询数据
	resp, err := svcCtx.UserModel.FindOneByUsername(ctx, user.Username)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
