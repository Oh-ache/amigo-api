package logic

import (
	"context"

	"amigo-api/app/sdk/rpc/internal/svc"
	"amigo-api/common/pb"
	osContext "amigo-api/common/utils/plug/objectsave/context"
	"amigo-api/common/utils/plug/objectsave/factory"
	"amigo-api/common/utils/plug/objectsave/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUploadTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUploadTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUploadTokenLogic {
	return &GetUploadTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUploadTokenLogic) GetUploadToken(in *pb.GetUploadTokenReq) (*pb.GetUploadTokenResp, error) {
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

	result, err := storageCtx.GetUploadToken(in.FileName)
	if err != nil {
		return nil, err
	}

	return &pb.GetUploadTokenResp{
		Token:  result.Token,
		Expire: 600,
	}, nil
}
