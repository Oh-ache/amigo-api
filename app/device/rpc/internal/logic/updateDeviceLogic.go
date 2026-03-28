package logic

import (
	"amigo-api/app/device/model"
	"context"

	"amigo-api/app/device/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	jsoniter "github.com/json-iterator/go"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/trace"
	"go.opentelemetry.io/otel/attribute"
)

type UpdateDeviceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateDeviceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDeviceLogic {
	return &UpdateDeviceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateDeviceLogic) UpdateDevice(in *pb.DeviceResp) (*pb.DeviceResp, error) {
	tracer := trace.TracerFromContext(l.ctx)
	ctx, span := tracer.Start(l.ctx, "开始更新设备")

	fast := jsoniter.ConfigFastest
	bytes2, _ := fast.Marshal(in)
	span.SetAttributes(
		attribute.String("update.device.param", string(bytes2)),
	)
	defer span.End()

	_, err := l.svcCtx.DeviceModel.FindOne(ctx, in.DeviceId)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, model.ErrNotFound
		}
		l.Errorf("Failed to find device by id %d: %v", in.DeviceId, err)
		return nil, err
	}

	var m model.Device
	if err := copier.Copy(&m, in); err != nil {
		l.Errorf("Failed to copy request data to model: %v", err)
		return nil, err
	}

	isDuplicate, err := l.svcCtx.DeviceModel.CheckDuplicate(ctx, &m)
	if err != nil {
		l.Errorf("Failed to check duplicate: %v", err)
		return nil, err
	}
	if isDuplicate {
		return nil, model.ErrDuplicate
	}

	if err := l.svcCtx.DeviceModel.Update(ctx, &m); err != nil {
		l.Errorf("Failed to update device: %v", err)
		return nil, err
	}

	var resp pb.DeviceResp
	if err := copier.Copy(&resp, &m); err != nil {
		l.Errorf("Failed to copy model to response: %v", err)
		return nil, err
	}

	span.SetAttributes(
		attribute.String("update.device.success", "ok"),
	)

	return &resp, nil
}
