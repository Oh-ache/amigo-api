package logic

import (
	"amigo-api/app/baseCode/model"
	"context"

	"amigo-api/app/baseCode/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	jsoniter "github.com/json-iterator/go"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/trace"
	"go.opentelemetry.io/otel/attribute"
)

type ListBaseCodeSortLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListBaseCodeSortLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListBaseCodeSortLogic {
	return &ListBaseCodeSortLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListBaseCodeSortLogic) ListBaseCodeSort(in *pb.ListBaseCodeSortReq) (*pb.ListBaseCodeSortResp, error) {
	// 从上下文中获取tracer
	tracer := trace.TracerFromContext(l.ctx)
	// 创建自定义span
	ctx, span := tracer.Start(l.ctx, "开始列表查询")
	// 设置span属性

	fast := jsoniter.ConfigFastest
	bytes2, _ := fast.Marshal(in)
	span.SetAttributes(
		attribute.String("listSort.param", string(bytes2)),
	)
	defer span.End()

	// 构造查询条件
	search := &model.BaseCodeSortSearch{}
	_ = copier.Copy(search, in)

	// 调用模型方法查询数据
	list, total, err := l.svcCtx.BaseCodeSortModel.List(ctx, search)
	if err != nil {
		l.Errorf("Failed to list base code sorts: %v", err)
		return nil, err
	}

	// 转换模型数据到响应格式
	var resp pb.ListBaseCodeSortResp
	resp.Total = total
	resp.List = make([]*pb.BaseCodeSortResp, 0, len(list))
	for _, sort := range list {
		var sortResp pb.BaseCodeSortResp
		if err := copier.Copy(&sortResp, sort); err != nil {
			l.Errorf("Failed to copy model to response: %v", err)
			continue
		}
		resp.List = append(resp.List, &sortResp)
	}

	span.SetAttributes(
		attribute.String("listSort.success", "ok"),
	)

	return &resp, nil
}
