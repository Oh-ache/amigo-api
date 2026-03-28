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

type UpdateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserLogic) UpdateUser(in *pb.UpdateUserReq) (*pb.SuccessResp, error) {
	// 从上下文中获取tracer
	tracer := trace.TracerFromContext(l.ctx)
	// 创建自定义span
	ctx, span := tracer.Start(l.ctx, "开始更新")
	// 设置span属性

	fast := jsoniter.ConfigFastest
	bytes2, _ := fast.Marshal(in)
	span.SetAttributes(
		attribute.String("update.param", string(bytes2)),
	)
	defer span.End()

	// 检查数据是否存在
	_, err := l.svcCtx.UserModel.FindOne(ctx, in.UserId)
	if err != nil {
		if err == model.ErrNotFound {
			return &pb.SuccessResp{Success: false}, model.ErrNotFound
		}
		l.Errorf("Failed to find user by id %d: %v", in.UserId, err)
		return &pb.SuccessResp{Success: false}, err
	}

	// 创建数据模型
	var m model.User
	if err := copier.Copy(&m, in); err != nil {
		l.Errorf("Failed to copy request data to model: %v", err)
		return &pb.SuccessResp{Success: false}, err
	}

	// 检查重复
	isDuplicate, err := l.svcCtx.UserModel.CheckDuplicate(ctx, &m)
	if err != nil {
		l.Errorf("Failed to check duplicate: %v", err)
		return &pb.SuccessResp{Success: false}, err
	}
	if isDuplicate {
		return &pb.SuccessResp{Success: false}, model.ErrDuplicate
	}

	// 更新数据
	if err := l.svcCtx.UserModel.Update(ctx, &m); err != nil {
		l.Errorf("Failed to update user: %v", err)
		return &pb.SuccessResp{Success: false}, err
	}

	span.SetAttributes(
		attribute.String("update.success", "ok"),
	)

	return &pb.SuccessResp{Success: true}, nil
}
