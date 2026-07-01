package logic

import (
	"context"
	"fmt"

	"amigo-api/app/baseCode/rpc/basecode"
	"amigo-api/app/sdk/rpc/internal/svc"
	"amigo-api/common/pb"
	osContext "amigo-api/common/utils/plug/objectsave/context"

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
	ossClient, err := l.svcCtx.GetOssClient()
	if err != nil {
		return nil, err
	}

	storageCtx := osContext.NewStorageContext(ossClient)

	_, err = storageCtx.UploadUrl(in.FileName, in.Url)
	if err != nil {
		return nil, err
	}

	host := GetBaseCode(l.ctx, l.svcCtx.BaseCodeRpc, "sdk", "ali.oss.host")

	return &pb.UploadUrlResp{
		Url: fmt.Sprintf("%s/%s", host, in.FileName),
	}, nil
}

func getBaseCode(ctx context.Context, baseCode basecode.BaseCode, sortKey, key string) string {
	item, _ := baseCode.GetBaseCode(ctx, &pb.GetBaseCodeReq{
		SortKey: sortKey,
		Key:     key,
	})
	if item != nil {
		return item.Content
	}
	return ""
}