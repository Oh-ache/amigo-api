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

type GetBaseCodeSortLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetBaseCodeSortLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBaseCodeSortLogic {
	return &GetBaseCodeSortLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetBaseCodeSortLogic) GetBaseCodeSort(in *pb.GetBaseCodeSortReq) (*pb.BaseCodeSortResp, error) {
	// 从上下文中获取tracer
	tracer := trace.TracerFromContext(l.ctx)
	// 创建自定义span
	ctx, span := tracer.Start(l.ctx, "开始查询")
	// 设置span属性

	fast := jsoniter.ConfigFastest
	bytes2, _ := fast.Marshal(in)
	span.SetAttributes(
		attribute.String("getSort.param", string(bytes2)),
	)
	defer span.End()

	var sort *model.BaseCodeSort
	var err error

	// 优先根据主键id查询
	if in.BaseCodeSortId != 0 {
		if sort, err = l.svcCtx.BaseCodeSortModel.FindOne(ctx, in.BaseCodeSortId); err == nil {
			return l.buildResponse(ctx, sort)
		} else if err != model.ErrNotFound {
			l.Errorf("Failed to find BaseCodeSort by id: %v", err)
			return nil, err
		}
		l.Infof("BaseCodeSort with id %d not found, checking by sort_key", in.BaseCodeSortId)
	}

	// 其次根据sort_key查询
	if in.SortKey != "" {
		if sort, err = l.svcCtx.BaseCodeSortModel.FindOneBySortKey(ctx, in.SortKey); err == nil {
			return l.buildResponse(ctx, sort)
		} else if err != model.ErrNotFound {
			l.Errorf("Failed to find BaseCodeSort by sort_key: %v", err)
			return nil, err
		}
		return nil, model.ErrNotFound
	}

	// 未提供查询条件
	return nil, model.ErrNotFound
}

// 统一构建响应对象
func (l *GetBaseCodeSortLogic) buildResponse(ctx context.Context, sort *model.BaseCodeSort) (*pb.BaseCodeSortResp, error) {
	var resp pb.BaseCodeSortResp
	if err := copier.Copy(&resp, sort); err != nil {
		l.Errorf("Failed to copy BaseCodeSort to BaseCodeSortResp: %v", err)
		return nil, err
	}
	return &resp, nil
}
