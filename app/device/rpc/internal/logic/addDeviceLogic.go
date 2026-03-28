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

type AddDeviceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddDeviceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddDeviceLogic {
	return &AddDeviceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddDeviceLogic) AddDevice(in *pb.AddDeviceReq) (*pb.DeviceResp, error) {
	tracer := trace.TracerFromContext(l.ctx)
	ctx, span := tracer.Start(l.ctx, "开始添加设备")

	fast := jsoniter.ConfigFastest
	bytes2, _ := fast.Marshal(in)
	span.SetAttributes(
		attribute.String("add.device.param", string(bytes2)),
	)
	defer span.End()

	var m model.Device
	if err := copier.Copy(&m, in); err != nil {
		l.Errorf("Failed to copy request data to model: %v", err)
		return nil, err
	}

	if m.IsDelete == 0 {
		m.IsDelete = 2
	}

	isDuplicate, err := l.svcCtx.DeviceModel.CheckDuplicate(ctx, &m)
	if err != nil {
		l.Errorf("Failed to check duplicate: %v", err)
		return nil, err
	}
	if isDuplicate {
		return nil, model.ErrDuplicate
	}

	result, err := l.svcCtx.DeviceModel.Insert(ctx, &m)
	if err != nil {
		l.Errorf("Failed to insert device: %v", err)
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		l.Errorf("Failed to get last insert id: %v", err)
		return nil, err
	}
	m.DeviceId = uint64(id)

	var resp pb.DeviceResp
	if err := copier.Copy(&resp, &m); err != nil {
		l.Errorf("Failed to copy model to response: %v", err)
		return nil, err
	}

	span.SetAttributes(
		attribute.String("add.device.success", "ok"),
	)

	return &resp, nil
}
