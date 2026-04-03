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

type ListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListLogic) List(req *types.ListBaseCodeItemReq) (resp *types.ListBaseCodeItemResp, err error) {
	// 从上下文中获取tracer
	tracer := trace.TracerFromContext(l.ctx)
	// 创建自定义span
	ctx, span := tracer.Start(l.ctx, "开始列表查询")
	// 设置span属性

	fast := jsoniter.ConfigFastest
	bytes2, _ := fast.Marshal(req)
	span.SetAttributes(
		attribute.String("list.param", string(bytes2)),
	)
	defer span.End()

	resp = &types.ListBaseCodeItemResp{}
	param := &pb.ListBaseCodeItemReq{}

	copier.Copy(param, req)
	rpcResp, err := l.svcCtx.BaseCodeRpcClient.ListBaseCodeItem(ctx, param)
	if err != nil {
		return nil, err
	}

	copier.Copy(resp, rpcResp)
	
	span.SetAttributes(
		attribute.String("list.success", "ok"),
	)
	
	return resp, nil
}
