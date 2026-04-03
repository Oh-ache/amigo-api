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

type GetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLogic {
	return &GetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLogic) Get(req *types.GetBaseCodeItemReq) (resp *types.GetBaseCodeItemResp, err error) {
	// 从上下文中获取tracer
	tracer := trace.TracerFromContext(l.ctx)
	// 创建自定义span
	ctx, span := tracer.Start(l.ctx, "开始查询")
	// 设置span属性

	fast := jsoniter.ConfigFastest
	bytes2, _ := fast.Marshal(req)
	span.SetAttributes(
		attribute.String("get.param", string(bytes2)),
	)
	defer span.End()

	resp = &types.GetBaseCodeItemResp{}
	param := &pb.GetBaseCodeItemReq{}

	copier.Copy(param, req)
	rpcResp, err := l.svcCtx.BaseCodeRpcClient.GetBaseCodeItem(ctx, param)
	if err != nil {
		return nil, err
	}

	copier.Copy(resp, rpcResp)
	
	span.SetAttributes(
		attribute.String("get.success", "ok"),
	)
	
	return resp, nil
}
