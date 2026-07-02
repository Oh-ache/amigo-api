package logic

import (
	"context"

	"amigo-api/app/device/model"
	"amigo-api/app/device/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListAppLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAppLogic {
	return &ListAppLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListAppLogic) ListApp(in *pb.ListAppReq) (*pb.ListAppResp, error) {
	search := &model.AppSearch{}
	if err := copier.Copy(search, in); err != nil {
		return nil, err
	}
	if search.PageSize <= 0 {
		search.PageSize = 10
	}
	if search.PageSize > 1000 {
		search.PageSize = 1000
	}

	list, total, err := l.svcCtx.AppModel.List(l.ctx, search)
	if err != nil {
		return nil, err
	}

	var respList []*pb.AppResp
	for _, item := range list {
		var respItem pb.AppResp
		if err := copier.Copy(&respItem, item); err != nil {
			continue
		}
		respList = append(respList, &respItem)
	}

	return &pb.ListAppResp{
		List:  respList,
		Total: total,
	}, nil
}
