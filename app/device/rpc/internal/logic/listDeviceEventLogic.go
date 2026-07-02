package logic

import (
	"context"

	"amigo-api/app/device/model"
	"amigo-api/app/device/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListDeviceEventLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListDeviceEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDeviceEventLogic {
	return &ListDeviceEventLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListDeviceEventLogic) ListDeviceEvent(in *pb.ListDeviceEventReq) (*pb.ListDeviceEventResp, error) {
	search := &model.DeviceEventSearch{
		DeviceId:   in.DeviceId,
		EventType:  in.EventType,
		EventLevel: in.EventLevel,
		Source:     in.Source,
		IsDelete:   in.IsDelete,
		Page:       in.Page,
		PageSize:   in.PageSize,
	}
	if search.PageSize <= 0 {
		search.PageSize = 10
	}
	if search.PageSize > 1000 {
		search.PageSize = 1000
	}
	if search.IsDelete == 0 {
		search.IsDelete = 2
	}

	list, total, err := l.svcCtx.DeviceEventModel.List(l.ctx, search)
	if err != nil {
		return nil, err
	}

	respList := make([]*pb.DeviceEventResp, 0, len(list))
	for _, item := range list {
		respList = append(respList, toPbDeviceEvent(item))
	}
	return &pb.ListDeviceEventResp{
		List:  respList,
		Total: total,
	}, nil
}
