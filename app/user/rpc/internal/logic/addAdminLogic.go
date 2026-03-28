package logic

import (
	"context"

	"amigo-api/app/user/model"
	"amigo-api/app/user/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	jsoniter "github.com/json-iterator/go"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/trace"
	"go.opentelemetry.io/otel/attribute"
)

type AddAdminLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddAdminLogic {
	return &AddAdminLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddAdminLogic) AddAdmin(in *pb.AddAdminReq) (*pb.SuccessResp, error) {
	// 从上下文中获取tracer
	tracer := trace.TracerFromContext(l.ctx)
	// 创建自定义span
	ctx, span := tracer.Start(l.ctx, "开始添加")
	// 设置span属性

	fast := jsoniter.ConfigFastest
	bytes2, _ := fast.Marshal(in)
	span.SetAttributes(
		attribute.String("add.param", string(bytes2)),
	)
	defer span.End()

	// 创建数据模型
	var m model.Admin
	if err := copier.Copy(&m, in); err != nil {
		l.Errorf("Failed to copy request data to model: %v", err)
		return &pb.SuccessResp{Success: false}, err
	}

	// 设置默认值
	if m.IsDelete == 0 {
		m.IsDelete = 2 // 2表示未删除
	}

	// 检查重复
	isDuplicate, err := l.svcCtx.AdminModel.CheckDuplicate(ctx, &m)
	if err != nil {
		l.Errorf("Failed to check duplicate: %v", err)
		return &pb.SuccessResp{Success: false}, err
	}
	if isDuplicate {
		return &pb.SuccessResp{Success: false}, model.ErrDuplicate
	}

	// 插入数据
	result, err := l.svcCtx.AdminModel.Insert(ctx, &m)
	if err != nil {
		l.Errorf("Failed to insert admin: %v", err)
		return &pb.SuccessResp{Success: false}, err
	}

	// 获取插入的ID
	id, err := result.LastInsertId()
	if err != nil {
		l.Errorf("Failed to get last insert id: %v", err)
		return &pb.SuccessResp{Success: false}, err
	}
	m.AdminId = uint64(id)

	span.SetAttributes(
		attribute.String("add.success", "ok"),
	)

	return &pb.SuccessResp{Success: true}, nil
}
