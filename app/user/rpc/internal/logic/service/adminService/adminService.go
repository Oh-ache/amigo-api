package adminService

import (
	"context"
	"fmt"

	"amigo-api/app/user/model"
	"amigo-api/app/user/rpc/internal/svc"
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

func Inert(ctx context.Context, svcCtx *svc.ServiceContext, admin *model.Admin, nodata bool) (*model.Admin, error) {
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
