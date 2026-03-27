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

type UpdateBaseCodeItemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateBaseCodeItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateBaseCodeItemLogic {
	return &UpdateBaseCodeItemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateBaseCodeItemLogic) UpdateBaseCodeItem(in *pb.BaseCodeItemResp) (*pb.BaseCodeItemResp, error) {
	// 从上下文中获取tracer
	tracer := trace.TracerFromContext(l.ctx)
	// 创建自定义span
	ctx, span := tracer.Start(l.ctx, "开始更新")
	// 设置span属性

	fast := jsoniter.ConfigFastest
	bytes2, _ := fast.Marshal(in)
	span.SetAttributes(
		attribute.String("updateItem.param", string(bytes2)),
	)
	defer span.End()

	// 检查数据是否存在
	_, err := l.svcCtx.BaseCodeItemModel.FindOne(ctx, in.BaseCodeItemId)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, model.ErrNotFound
		}
		l.Errorf("Failed to find base code item by id %d: %v", in.BaseCodeItemId, err)
		return nil, err
	}

	// 创建数据模型
	var m model.BaseCodeItem
	if err := copier.Copy(&m, in); err != nil {
		l.Errorf("Failed to copy request data to model: %v", err)
		return nil, err
	}

	// 检查重复
	isDuplicate, err := l.svcCtx.BaseCodeItemModel.CheckDuplicate(ctx, &m)
	if err != nil {
		l.Errorf("Failed to check duplicate: %v", err)
		return nil, err
	}
	if isDuplicate {
		return nil, model.ErrDuplicate
	}

	// 更新数据
	if err := l.svcCtx.BaseCodeItemModel.Update(ctx, &m); err != nil {
		l.Errorf("Failed to update base code item: %v", err)
		return nil, err
	}

	// 构造响应
	var resp pb.BaseCodeItemResp
	if err := copier.Copy(&resp, &m); err != nil {
		l.Errorf("Failed to copy model to response: %v", err)
		return nil, err
	}

	span.SetAttributes(
		attribute.String("updateItem.success", "ok"),
	)

	return &resp, nil
}
