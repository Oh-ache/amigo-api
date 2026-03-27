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

type ListBaseCodeItemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListBaseCodeItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListBaseCodeItemLogic {
	return &ListBaseCodeItemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListBaseCodeItemLogic) ListBaseCodeItem(in *pb.ListBaseCodeItemReq) (*pb.ListBaseCodeItemResp, error) {
	// 从上下文中获取tracer
	tracer := trace.TracerFromContext(l.ctx)
	// 创建自定义span
	ctx, span := tracer.Start(l.ctx, "开始列表查询")
	// 设置span属性

	fast := jsoniter.ConfigFastest
	bytes2, _ := fast.Marshal(in)
	span.SetAttributes(
		attribute.String("listItem.param", string(bytes2)),
	)
	defer span.End()

	// 构造查询条件
	search := &model.BaseCodeItemSearch{}
	_ = copier.Copy(search, in)

	// 调用模型方法查询数据
	list, total, err := l.svcCtx.BaseCodeItemModel.List(ctx, search)
	if err != nil {
		l.Errorf("Failed to list base code items: %v", err)
		return nil, err
	}

	// 转换模型数据到响应格式
	var resp pb.ListBaseCodeItemResp
	resp.Total = total
	resp.List = make([]*pb.BaseCodeItemResp, 0, len(list))
	for _, item := range list {
		var itemResp pb.BaseCodeItemResp
		if err := copier.Copy(&itemResp, item); err != nil {
			l.Errorf("Failed to copy model to response: %v", err)
			continue
		}
		resp.List = append(resp.List, &itemResp)
	}

	span.SetAttributes(
		attribute.String("listItem.success", "ok"),
	)

	return &resp, nil
}
