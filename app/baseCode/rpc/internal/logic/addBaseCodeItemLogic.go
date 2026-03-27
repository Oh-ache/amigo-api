package logic

import (
	"context"

	"amigo-api/app/baseCode/model"

	"amigo-api/app/baseCode/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	jsoniter "github.com/json-iterator/go"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/trace"
	"go.opentelemetry.io/otel/attribute"
)

type AddBaseCodeItemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddBaseCodeItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddBaseCodeItemLogic {
	return &AddBaseCodeItemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddBaseCodeItemLogic) AddBaseCodeItem(in *pb.AddBaseCodeItemReq) (*pb.BaseCodeItemResp, error) {
	// 从上下文中获取tracer
	tracer := trace.TracerFromContext(l.ctx)

	// 创建自定义span
	ctx, span := tracer.Start(l.ctx, "开始添加")
	// 设置span属性

	fast := jsoniter.ConfigFastest
	bytes2, _ := fast.Marshal(in)
	span.SetAttributes(
		attribute.String("addItem.param", string(bytes2)),
	)
	defer span.End()

	var m model.BaseCodeItem
	if err := copier.Copy(&m, in); err != nil {
		l.Errorf("Failed to copy request data to model: %v", err)
		return nil, err
	}

	if _, err := l.svcCtx.BaseCodeItemModel.Insert(ctx, &m); err != nil {
		l.Errorf("Failed to insert base code item: %v", err)
		return nil, err
	}

	var resp pb.BaseCodeItemResp
	if err := copier.Copy(&resp, &m); err != nil {
		l.Errorf("Failed to copy model to response: %v", err)
		return nil, err
	}

	span.SetAttributes(
		attribute.String("addItem.success", "ok"),
	)

	return &resp, nil
}
