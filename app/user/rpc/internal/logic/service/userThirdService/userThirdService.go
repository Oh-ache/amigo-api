package userThirdService

import (
	"context"
	"encoding/json"
	"fmt"

	"amigo-api/app/user/model"
	"amigo-api/app/user/rpc/internal/logic/service/userService"
	"amigo-api/app/user/rpc/internal/svc"
	"amigo-api/common/pb"
	"amigo-api/common/utils"
	"amigo-api/common/utils/plug/miniapp"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
)

func Code2Session(ctx context.Context, svcCtx *svc.ServiceContext, appType, code string) (*pb.MiniappCodeResp, error) {
	// FIXME 考虑加缓存或者幂等处理
	client := &miniapp.Client{}
	if client, _ = miniapp.GetClient(appType); client == nil {
		conf, _ := svcCtx.BaseCodeRpc.GetBaseCode(ctx, &pb.GetBaseCodeReq{
			SortKey: "user_third",
			Key:     appType,
		})

		miniappConf := &miniapp.MiniAppConfig{}
		json.Unmarshal([]byte(conf.Content), miniappConf)

		client, _ = miniapp.InitClient(appType, miniappConf)
	}

	session, _ := client.Auth.Session(ctx, code)

	resp := &pb.MiniappCodeResp{}
	_ = copier.Copy(resp, session)
	return resp, nil
}

func InsertUser(ctx context.Context, svcCtx *svc.ServiceContext, appType, openid, unionid string) (*model.User, error) {
	userThirdParty, err := svcCtx.UserThirdParty.FindOneByAppTypeOpenid(ctx, appType, openid)
	if err != nil && err != sqlc.ErrNotFound {
		return nil, err
	}

	if userThirdParty == nil {
		// FIXME 可能要去重
		username := fmt.Sprintf("%s%s", "用户", utils.GetRandomString(6))
		insert := &model.User{
			Username: username,
			Avatar:   "http://dummyimage.com/100x100",
		}

		user := &model.User{}
		if user, err = userService.Insert(ctx, svcCtx, insert, false); err != nil {
			return nil, err
		}

		baseCode, _ := svcCtx.BaseCodeRpc.GetBaseCode(ctx, &pb.GetBaseCodeReq{
			SortKey: "user_third",
			Key:     appType,
		})
		userThirdParty := &model.UserThirdParty{
			UserId:       user.UserId,
			Openid:       openid,
			Unionid:      unionid,
			Platform:     baseCode.Content1,
			PlatformType: baseCode.Content2,
			AppType:      appType,
		}
		if _, err = svcCtx.UserThirdParty.Insert(ctx, userThirdParty); err != nil {
			return nil, err
		}
		return user, nil
	}

	user, _ := svcCtx.UserModel.FindOne(ctx, userThirdParty.UserId)
	if user.IsDelete == 1 {
		return nil, fmt.Errorf("用户已注销")
	}
	return user, nil
}
