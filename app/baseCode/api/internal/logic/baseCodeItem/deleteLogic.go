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

type DeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLogic) Delete(req *types.DeleteBaseCodeItemReq) (resp *types.EmptyResp, err error) {
	// 从上下文中获取tracer
	tracer := trace.TracerFromContext(l.ctx)
	// 创建自定义span
	ctx, span := tracer.Start(l.ctx, "开始删除")
	// 设置span属性

	fast := jsoniter.ConfigFastest
	bytes2, _ := fast.Marshal(req)
	span.SetAttributes(
		attribute.String("delete.param", string(bytes2)),
	)
	defer span.End()

	resp = &types.EmptyResp{}
	param := &pb.DeleteBaseCodeItemReq{}

	copier.Copy(param, req)
	if _, err := l.svcCtx.BaseCodeRpcClient.DeleteBaseCodeItem(ctx, param); err != nil {
		return nil, err
	}

	span.SetAttributes(
		attribute.String("delete.success", "ok"),
	)

	return resp, nil
}
