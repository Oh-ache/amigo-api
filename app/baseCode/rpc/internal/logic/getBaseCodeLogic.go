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

type GetBaseCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetBaseCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBaseCodeLogic {
	return &GetBaseCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetBaseCodeLogic) GetBaseCode(in *pb.GetBaseCodeReq) (*pb.BaseCodeResp, error) {
	// 从上下文中获取tracer
	tracer := trace.TracerFromContext(l.ctx)
	// 创建自定义span
	ctx, span := tracer.Start(l.ctx, "开始查询")
	// 设置span属性

	fast := jsoniter.ConfigFastest
	bytes2, _ := fast.Marshal(in)
	span.SetAttributes(
		attribute.String("get.param", string(bytes2)),
	)
	defer span.End()

	// 根据主键id查询数据
	if in.BaseCodeId != 0 {
		data, err := l.svcCtx.BaseCodeModel.FindOne(ctx, in.BaseCodeId)
		if err != nil {
			if err == model.ErrNotFound {
				return nil, model.ErrNotFound
			}
			l.Errorf("Failed to find base code by id %d: %v", in.BaseCodeId, err)
			return nil, err
		}
		var resp pb.BaseCodeResp
		if err := copier.Copy(&resp, data); err != nil {
			l.Errorf("Failed to copy model to response: %v", err)
			return nil, err
		}
		// 设置span属性
		span.SetAttributes(
			attribute.String("get.success", "ok"),
		)
		return &resp, nil
	}

	// 根据sort_key和key查询数据
	if in.SortKey != "" && in.Key != "" {
		data, err := l.svcCtx.BaseCodeModel.FindOneBySortKeyKey(ctx, in.SortKey, in.Key)
		if err != nil {
			if err == model.ErrNotFound {
				return nil, model.ErrNotFound
			}
			l.Errorf("Failed to find base code by sort_key %s and key %s: %v", in.SortKey, in.Key, err)
			return nil, err
		}
		var resp pb.BaseCodeResp
		if err := copier.Copy(&resp, data); err != nil {
			l.Errorf("Failed to copy model to response: %v", err)
			return nil, err
		}
		// 设置span属性
		span.SetAttributes(
			attribute.String("get.success", "ok"),
		)
		return &resp, nil
	}

	// 参数无效
	l.Errorf("Invalid parameters: base_code_id, sort_key or key must be provided")
	return nil, model.ErrNotFound
}
