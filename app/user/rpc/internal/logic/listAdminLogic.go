package logic

import (
	"context"

	"amigo-api/app/user/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListAdminLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAdminLogic {
	return &ListAdminLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListAdminLogic) ListAdmin(in *pb.ListAdminReq) (*pb.ListAdminResp, error) {
	// todo: add your logic here and delete this line

	return &pb.ListAdminResp{}, nil
}
