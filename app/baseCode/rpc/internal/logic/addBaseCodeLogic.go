package logic

import (
	"context"

	"amigo-api/app/baseCode/model"
	"amigo-api/app/baseCode/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	jsoniter "github.com/json-iterator/go"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/trace"
	"go.opentelemetry.io/otel/attribute"
)

type AddBaseCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddBaseCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddBaseCodeLogic {
	return &AddBaseCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddBaseCodeLogic) AddBaseCode(in *pb.AddBaseCodeReq) (*pb.BaseCodeResp, error) {
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
	var m model.BaseCode
	if err := copier.Copy(&m, in); err != nil {
		l.Errorf("Failed to copy request data to model: %v", err)
		return nil, err
	}

	// 设置默认值
	if m.IsDelete == 0 {
		m.IsDelete = 2 // 2表示未删除
	}

	// 检查重复
	isDuplicate, err := l.svcCtx.BaseCodeModel.CheckDuplicate(ctx, &m)
	if err != nil {
		l.Errorf("Failed to check duplicate: %v", err)
		return nil, err
	}
	if isDuplicate {
		return nil, model.ErrDuplicate
	}

	// 插入数据
	result, err := l.svcCtx.BaseCodeModel.Insert(ctx, &m)
	if err != nil {
		l.Errorf("Failed to insert base code: %v", err)
		return nil, err
	}

	// 获取插入的ID
	id, err := result.LastInsertId()
	if err != nil {
		l.Errorf("Failed to get last insert id: %v", err)
		return nil, err
	}
	m.BaseCodeId = uint64(id)

	// 构造响应
	var resp pb.BaseCodeResp
	if err := copier.Copy(&resp, &m); err != nil {
		l.Errorf("Failed to copy model to response: %v", err)
		return nil, err
	}

	span.SetAttributes(
		attribute.String("add.success", "ok"),
	)

	return &resp, nil
}
