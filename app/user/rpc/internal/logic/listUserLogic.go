package logic

import (
	"context"

	"amigo-api/app/user/rpc/internal/svc"
	"amigo-api/app/user/model"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	jsoniter "github.com/json-iterator/go"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/trace"
	"go.opentelemetry.io/otel/attribute"
)

type ListUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUserLogic {
	return &ListUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListUserLogic) ListUser(in *pb.ListUserReq) (*pb.ListUserResp, error) {
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
	search := &model.UserSearch{}
	if err := copier.Copy(search, in); err != nil {
		l.Errorf("Failed to copy request to search: %v", err)
		return nil, err
	}

	// 查询数据
	list, total, err := l.svcCtx.UserModel.List(ctx, search)
	if err != nil {
		l.Errorf("Failed to list users: %v", err)
		return nil, err
	}

	// 构造响应
	resp := &pb.ListUserResp{}
	if err := copier.Copy(resp, &struct {
		List  []*model.User
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
