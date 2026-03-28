package logic

import (
	"context"
	"fmt"

	"amigo-api/app/baseCode/rpc/basecode"
	"amigo-api/app/sdk/rpc/internal/svc"
	"amigo-api/common/pb"
	osContext "amigo-api/common/utils/plug/objectsave/context"
	"amigo-api/common/utils/plug/objectsave/factory"
	"amigo-api/common/utils/plug/objectsave/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadUrlLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUploadUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadUrlLogic {
	return &UploadUrlLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UploadUrlLogic) UploadUrl(in *pb.UploadUrlReq) (*pb.UploadUrlResp, error) {
	factory := factory.NewStorageFactory()

	storageConfig := &model.StorageConfig{
		Type: "oss",
		OssConfig: &model.OssConfig{
			Endpoint:        GetBaseCode(l.ctx, l.svcCtx.BaseCodeRpc, "sdk", "ali.oss.endpoint"),
			AccessKeyId:     GetBaseCode(l.ctx, l.svcCtx.BaseCodeRpc, "sdk", "ali.accessKey"),
			AccessKeySecret: GetBaseCode(l.ctx, l.svcCtx.BaseCodeRpc, "sdk", "ali.accessKeySecret"),
			Bucket:          GetBaseCode(l.ctx, l.svcCtx.BaseCodeRpc, "sdk", "ali.oss.bucket"),
			Region:          GetBaseCode(l.ctx, l.svcCtx.BaseCodeRpc, "sdk", "ali.oss.region"),
		},
	}

	ossClient, _ := factory.CreateClient(storageConfig)
	storageCtx := osContext.NewStorageContext(ossClient)

	_, err := storageCtx.UploadUrl(in.FileName, in.Url)
	if err != nil {
		return nil, err
	}

	host := GetBaseCode(l.ctx, l.svcCtx.BaseCodeRpc, "sdk", "ali.oss.host")

	return &pb.UploadUrlResp{
		Url: fmt.Sprintf("%s/%s", host, in.FileName),
	}, nil
}

// TODO 后面添加到公共方法
func GetBaseCode(ctx context.Context, baseCode basecode.BaseCode, sortKey, key string) string {
	getBaseCodeReq := &pb.GetBaseCodeReq{}
	getBaseCodeReq.SortKey = sortKey
	getBaseCodeReq.Key = key

	item, _ := baseCode.GetBaseCode(ctx, getBaseCodeReq)
	return item.Content
}
