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

type GetDeviceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDeviceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDeviceLogic {
	return &GetDeviceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetDeviceLogic) GetDevice(in *pb.GetDeviceReq) (*pb.DeviceResp, error) {
	tracer := trace.TracerFromContext(l.ctx)
	ctx, span := tracer.Start(l.ctx, "开始获取设备")

	fast := jsoniter.ConfigFastest
	bytes2, _ := fast.Marshal(in)
	span.SetAttributes(
		attribute.String("get.device.param", string(bytes2)),
	)
	defer span.End()

	var device *model.Device
	var err error

	if in.DeviceId != 0 {
		device, err = l.svcCtx.DeviceModel.FindOne(ctx, in.DeviceId)
	} else if in.MacAddress != "" {
		device, err = l.svcCtx.DeviceModel.FindOneByMacAddress(ctx, in.MacAddress)
	} else {
		return nil, model.ErrNotFound
	}

	if err != nil {
		if err == model.ErrNotFound {
			return nil, model.ErrNotFound
		}
		l.Errorf("Failed to get device: %v", err)
		return nil, err
	}

	var resp pb.DeviceResp
	if err := copier.Copy(&resp, device); err != nil {
		l.Errorf("Failed to copy device: %v", err)
		return nil, err
	}

	span.SetAttributes(
		attribute.String("get.device.success", "ok"),
	)

	return &resp, nil
}
