package logic

import (
	"context"
	"fmt"

	"amigo-api/app/sdk/rpc/internal/svc"
	"amigo-api/common/pb"
	osContext "amigo-api/common/utils/plug/objectsave/context"
	"amigo-api/common/utils/plug/objectsave/factory"
	"amigo-api/common/utils/plug/objectsave/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUploadFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadFileLogic {
	return &UploadFileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UploadFileLogic) UploadFile(in *pb.UploadFileReq) (*pb.UploadFileResp, error) {
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

	host := GetBaseCode(l.ctx, l.svcCtx.BaseCodeRpc, "sdk", "ali.oss.host")

	_, err := storageCtx.UploadFile(in.FileName, in.File)
	if err != nil {
		return nil, err
	}

	return &pb.UploadFileResp{
		Url: fmt.Sprintf("%s/%s", host, in.FileName),
	}, nil
}
