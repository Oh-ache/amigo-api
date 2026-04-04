package logic

import (
	"context"

	"amigo-api/app/user/model"
	"amigo-api/app/user/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
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
	// 创建数据模型
	var m model.Admin
	if err := copier.Copy(&m, in); err != nil {
		return &pb.SuccessResp{Success: false}, err
	}

	// 设置默认值
	if m.IsDelete == 0 {
		m.IsDelete = 2 // 2表示未删除
	}

	// 检查重复
	isDuplicate, err := l.svcCtx.AdminModel.CheckDuplicate(l.ctx, &m)
	if err != nil {
		return &pb.SuccessResp{Success: false}, err
	}
	if isDuplicate {
		return &pb.SuccessResp{Success: false}, model.ErrDuplicate
	}

	// 插入数据
	result, err := l.svcCtx.AdminModel.Insert(l.ctx, &m)
	if err != nil {
		return &pb.SuccessResp{Success: false}, err
	}

	// 获取插入的ID
	id, err := result.LastInsertId()
	if err != nil {
		return &pb.SuccessResp{Success: false}, err
	}
	m.AdminId = uint64(id)

	return &pb.SuccessResp{Success: true}, nil
}
