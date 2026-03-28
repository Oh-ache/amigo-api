package logic

import (
	"context"

	"amigo-api/app/device/model"
	"amigo-api/app/device/rpc/internal/svc"
	"amigo-api/common/pb"

	jsoniter "github.com/json-iterator/go"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/trace"
	"go.opentelemetry.io/otel/attribute"
)

type DeleteDeviceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteDeviceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDeviceLogic {
	return &DeleteDeviceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteDeviceLogic) DeleteDevice(in *pb.DeleteDeviceReq) (*pb.DeleteDeviceResp, error) {
	tracer := trace.TracerFromContext(l.ctx)
	ctx, span := tracer.Start(l.ctx, "开始删除设备")

	fast := jsoniter.ConfigFastest
	bytes2, _ := fast.Marshal(in)
	span.SetAttributes(
		attribute.String("delete.device.param", string(bytes2)),
	)
	defer span.End()

	var device *model.Device
	var err error

	if in.DeviceId != 0 {
		device, err = l.svcCtx.DeviceModel.FindOne(ctx, in.DeviceId)
	} else if in.MacAddress != "" {
		device, err = l.svcCtx.DeviceModel.FindOneByMacAddress(ctx, in.MacAddress)
	} else {
		return &pb.DeleteDeviceResp{Success: false}, nil
	}

	if err != nil {
		if err == model.ErrNotFound {
			return &pb.DeleteDeviceResp{Success: false}, nil
		}
		l.Errorf("Failed to find device for deletion: %v", err)
		return nil, err
	}

	if err := l.svcCtx.DeviceModel.Delete(ctx, device.DeviceId); err != nil {
		l.Errorf("Failed to delete device: %v", err)
		return nil, err
	}

	span.SetAttributes(
		attribute.String("delete.device.success", "ok"),
	)

	return &pb.DeleteDeviceResp{Success: true}, nil
}
