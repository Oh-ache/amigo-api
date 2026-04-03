package baseCodeItem

import (
	"context"

	"amigo-api/app/baseCode/api/internal/svc"
	"amigo-api/app/baseCode/api/internal/types"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	jsoniter "github.com/json-iterator/go"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/trace"
	"go.opentelemetry.io/otel/attribute"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.UpdateBaseCodeItemReq) (resp *types.EmptyResp, err error) {
	// 从上下文中获取tracer
	tracer := trace.TracerFromContext(l.ctx)
	// 创建自定义span
	ctx, span := tracer.Start(l.ctx, "开始更新")
	// 设置span属性

	fast := jsoniter.ConfigFastest
	bytes2, _ := fast.Marshal(req)
	span.SetAttributes(
		attribute.String("update.param", string(bytes2)),
	)
	defer span.End()

	resp = &types.EmptyResp{}
	// First, get the existing item to populate the BaseCodeItemResp
	getReq := &pb.GetBaseCodeItemReq{
		BaseCodeItemId: req.BaseCodeItemId,
	}
	existingItem, err := l.svcCtx.BaseCodeRpcClient.GetBaseCodeItem(ctx, getReq)
	if err != nil {
		return nil, err
	}

	// Copy the update request fields into the existing item
	copier.Copy(existingItem, req)

	if _, err := l.svcCtx.BaseCodeRpcClient.UpdateBaseCodeItem(ctx, existingItem); err != nil {
		return nil, err
	}

	span.SetAttributes(
		attribute.String("update.success", "ok"),
	)

	return resp, nil
}
