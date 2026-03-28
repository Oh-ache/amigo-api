package logic

import (
	"context"

	"amigo-api/app/user/rpc/internal/svc"
	"amigo-api/app/user/model"
	"amigo-api/common/pb"

	jsoniter "github.com/json-iterator/go"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/trace"
	"go.opentelemetry.io/otel/attribute"
)

type DeleteAdminLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAdminLogic {
	return &DeleteAdminLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteAdminLogic) DeleteAdmin(in *pb.DeleteAdminReq) (*pb.SuccessResp, error) {
	// 从上下文中获取tracer
	tracer := trace.TracerFromContext(l.ctx)
	// 创建自定义span
	ctx, span := tracer.Start(l.ctx, "开始删除")
	// 设置span属性

	fast := jsoniter.ConfigFastest
	bytes2, _ := fast.Marshal(in)
	span.SetAttributes(
		attribute.String("delete.param", string(bytes2)),
	)
	defer span.End()

	// 检查主键id是否存在
	if in.AdminId == 0 {
		return &pb.SuccessResp{Success: false}, model.ErrNotFound
	}

	// 根据主键id删除数据
	if err := l.svcCtx.AdminModel.Delete(ctx, in.AdminId); err != nil {
		if err == model.ErrNotFound {
			return &pb.SuccessResp{Success: false}, model.ErrNotFound
		}
		l.Errorf("Failed to delete Admin by id %d: %v", in.AdminId, err)
		return &pb.SuccessResp{Success: false}, err
	}

	// 删除成功
	span.SetAttributes(
		attribute.String("delete.success", "ok"),
	)

	return &pb.SuccessResp{Success: true}, nil
}
