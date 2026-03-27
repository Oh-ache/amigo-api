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

type GetBaseCodeItemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetBaseCodeItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBaseCodeItemLogic {
	return &GetBaseCodeItemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetBaseCodeItemLogic) GetBaseCodeItem(in *pb.GetBaseCodeItemReq) (*pb.BaseCodeItemResp, error) {
	// 从上下文中获取tracer
	tracer := trace.TracerFromContext(l.ctx)
	// 创建自定义span
	ctx, span := tracer.Start(l.ctx, "开始查询")
	// 设置span属性

	fast := jsoniter.ConfigFastest
	bytes2, _ := fast.Marshal(in)
	span.SetAttributes(
		attribute.String("getItem.param", string(bytes2)),
	)
	defer span.End()

	var item *model.BaseCodeItem
	var err error

	if in.BaseCodeItemId != 0 {
		item, err = l.svcCtx.BaseCodeItemModel.FindOne(ctx, in.BaseCodeItemId)
	} else if in.SortKey != "" && in.Key != "" {
		item, err = l.svcCtx.BaseCodeItemModel.FindOneBySortKeyKey(ctx, in.SortKey, in.Key)
	} else {
		return nil, model.ErrInvalidParams
	}

	if err != nil {
		if err == model.ErrNotFound {
			return nil, model.ErrNotFound
		}
		l.Errorf("Failed to get base code item: %v", err)
		return nil, err
	}

	var resp pb.BaseCodeItemResp
	_ = copier.Copy(&resp, item)

	span.SetAttributes(
		attribute.String("getItem.success", "ok"),
	)

	return &resp, nil
}
