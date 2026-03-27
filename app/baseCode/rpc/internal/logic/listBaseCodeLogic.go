package logic

import (
	"context"

	"amigo-api/app/baseCode/rpc/internal/svc"
	"amigo-api/app/baseCode/model"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	jsoniter "github.com/json-iterator/go"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/trace"
	"go.opentelemetry.io/otel/attribute"
)

type ListBaseCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListBaseCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListBaseCodeLogic {
	return &ListBaseCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListBaseCodeLogic) ListBaseCode(in *pb.ListBaseCodeReq) (*pb.ListBaseCodeResp, error) {
	// 从上下文中获取tracer
	tracer := trace.TracerFromContext(l.ctx)
	// 创建自定义span
	ctx, span := tracer.Start(l.ctx, "开始列表查询")
	// 设置span属性

	fast := jsoniter.ConfigFastest
	bytes2, _ := fast.Marshal(in)
	span.SetAttributes(
		attribute.String("list.param", string(bytes2)),
	)
	defer span.End()

	// 构建查询条件
	search := &model.BaseCodeSearch{}
	if err := copier.Copy(search, in); err != nil {
		l.Errorf("Failed to copy request to search: %v", err)
		return nil, err
	}

	// 查询数据
	list, total, err := l.svcCtx.BaseCodeModel.List(ctx, search)
	if err != nil {
		l.Errorf("Failed to list base codes: %v", err)
		return nil, err
	}

	// 构造响应
	resp := &pb.ListBaseCodeResp{}
	if err := copier.Copy(resp, &struct {
		List  []*model.BaseCode
		Total int64
	}{
		List:  list,
		Total: total,
	}); err != nil {
		l.Errorf("Failed to copy list to response: %v", err)
		return nil, err
	}

	span.SetAttributes(
		attribute.String("list.success", "ok"),
	)

	return resp, nil
}
