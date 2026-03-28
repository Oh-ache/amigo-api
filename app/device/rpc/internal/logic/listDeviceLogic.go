package logic

import (
	"context"

	"amigo-api/app/device/model"
	"amigo-api/app/device/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	jsoniter "github.com/json-iterator/go"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/trace"
	"go.opentelemetry.io/otel/attribute"
)

type ListDeviceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListDeviceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDeviceLogic {
	return &ListDeviceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListDeviceLogic) ListDevice(in *pb.ListDeviceReq) (*pb.ListDeviceResp, error) {
	tracer := trace.TracerFromContext(l.ctx)
	ctx, span := tracer.Start(l.ctx, "开始查询设备列表")

	fast := jsoniter.ConfigFastest
	bytes2, _ := fast.Marshal(in)
	span.SetAttributes(
		attribute.String("list.device.param", string(bytes2)),
	)
	defer span.End()

	search := &model.DeviceSearch{
		Name:       in.Name,
		UserId:     in.UserId,
		MacAddress: in.MacAddress,
		InternalIp: in.InternalIp,
		IsRunning:  in.IsRunning,
		IsDelete:   in.IsDelete,
		Page:       in.Page,
		PageSize:   in.PageSize,
	}

	list, total, err := l.svcCtx.DeviceModel.List(ctx, search)
	if err != nil {
		l.Errorf("Failed to list devices: %v", err)
		return nil, err
	}

	var respList []*pb.DeviceResp
	for _, item := range list {
		var respItem pb.DeviceResp
		if err := copier.Copy(&respItem, item); err != nil {
			l.Errorf("Failed to copy device item: %v", err)
			return nil, err
		}
		respList = append(respList, &respItem)
	}

	span.SetAttributes(
		attribute.String("list.device.success", "ok"),
	)

	return &pb.ListDeviceResp{
		List:  respList,
		Total: total,
	}, nil
}
