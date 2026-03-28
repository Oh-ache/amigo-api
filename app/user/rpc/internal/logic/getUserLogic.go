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

type GetUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserLogic) GetUser(in *pb.GetUserReq) (*pb.UserResp, error) {
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
	if in.UserId != 0 {
		data, err := l.svcCtx.UserModel.FindOne(ctx, in.UserId)
		if err != nil {
			if err == model.ErrNotFound {
				return nil, model.ErrNotFound
			}
			l.Errorf("Failed to find user by id %d: %v", in.UserId, err)
			return nil, err
		}
		var resp pb.UserResp
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
		data, err := l.svcCtx.UserModel.FindOneByUsername(ctx, in.Username)
		if err != nil {
			if err == model.ErrNotFound {
				return nil, model.ErrNotFound
			}
			l.Errorf("Failed to find user by username %s: %v", in.Username, err)
			return nil, err
		}
		var resp pb.UserResp
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
	l.Errorf("Invalid parameters: user_id or username must be provided")
	return nil, model.ErrNotFound
}
