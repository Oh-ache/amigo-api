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

type UpdateBaseCodeSortLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateBaseCodeSortLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateBaseCodeSortLogic {
	return &UpdateBaseCodeSortLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateBaseCodeSortLogic) UpdateBaseCodeSort(in *pb.BaseCodeSortResp) (*pb.BaseCodeSortResp, error) {
	// 从上下文中获取tracer
	tracer := trace.TracerFromContext(l.ctx)
	// 创建自定义span
	ctx, span := tracer.Start(l.ctx, "开始更新")
	// 设置span属性

	fast := jsoniter.ConfigFastest
	bytes2, _ := fast.Marshal(in)
	span.SetAttributes(
		attribute.String("updateSort.param", string(bytes2)),
	)
	defer span.End()

	// 创建 model.BaseCodeSort 实例
	var sort model.BaseCodeSort

	// 使用 copier 库将入参复制到 sort 实例
	if err := copier.Copy(&sort, in); err != nil {
		l.Errorf("Failed to copy BaseCodeSortResp to BaseCodeSort: %v", err)
		return nil, err
	}

	// 先判断是否重复
	isDuplicate, err := l.svcCtx.BaseCodeSortModel.CheckDuplicate(ctx, &sort)
	if err != nil {
		l.Errorf("Failed to check duplicate for BaseCodeSort: %v", err)
		return nil, err
	}
	if isDuplicate {
		return nil, model.ErrDuplicate
	}

	// 更新数据
	if err := l.svcCtx.BaseCodeSortModel.Update(ctx, &sort); err != nil {
		l.Errorf("Failed to update BaseCodeSort: %v", err)
		return nil, err
	}

	// 查询更新后的完整数据
	updatedSort, err := l.svcCtx.BaseCodeSortModel.FindOne(ctx, sort.BaseCodeSortId)
	if err != nil {
		l.Errorf("Failed to find updated BaseCodeSort: %v", err)
		return nil, err
	}

	// 创建返回响应
	var resp pb.BaseCodeSortResp
	if err := copier.Copy(&resp, updatedSort); err != nil {
		l.Errorf("Failed to copy BaseCodeSort to BaseCodeSortResp: %v", err)
		return nil, err
	}

	span.SetAttributes(
		attribute.String("updateSort.success", "ok"),
	)

	return &resp, nil
}
