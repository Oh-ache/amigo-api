package logic

import (
	"context"

	"amigo-api/app/sdk/rpc/internal/svc"
	"amigo-api/common/pb"
	osContext "amigo-api/common/utils/plug/objectsave/context"

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
	ossClient, err := l.svcCtx.GetOssClient()
	if err != nil {
		return nil, err
	}

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
