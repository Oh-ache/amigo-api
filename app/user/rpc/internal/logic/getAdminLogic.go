package logic

import (
	"amigo-api/app/user/model"
	"context"

	"amigo-api/app/user/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	jsoniter "github.com/json-iterator/go"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/trace"
	"go.opentelemetry.io/otel/attribute"
)

type GetAdminLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAdminLogic {
	return &GetAdminLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAdminLogic) GetAdmin(in *pb.GetAdminReq) (*pb.AdminResp, error) {
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
	if in.AdminId != 0 {
		data, err := l.svcCtx.AdminModel.FindOne(ctx, in.AdminId)
		if err != nil {
			if err == model.ErrNotFound {
				return nil, model.ErrNotFound
			}
			l.Errorf("Failed to find admin by id %d: %v", in.AdminId, err)
			return nil, err
		}
		var resp pb.AdminResp
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

	// 根据username查询数据
	if in.Username != "" {
		data, err := l.svcCtx.AdminModel.FindOneByUsername(ctx, in.Username)
		if err != nil {
			if err == model.ErrNotFound {
				return nil, model.ErrNotFound
			}
			l.Errorf("Failed to find admin by username %s: %v", in.Username, err)
			return nil, err
		}
		var resp pb.AdminResp
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

	// 根据mobile查询数据
	if in.Mobile != "" {
		data, err := l.svcCtx.AdminModel.FindOneByMobile(ctx, in.Mobile)
		if err != nil {
			if err == model.ErrNotFound {
				return nil, model.ErrNotFound
			}
			l.Errorf("Failed to find admin by mobile %s: %v", in.Mobile, err)
			return nil, err
		}
		var resp pb.AdminResp
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
	l.Errorf("Invalid parameters: admin_id, username or mobile must be provided")
	return nil, model.ErrNotFound
}
