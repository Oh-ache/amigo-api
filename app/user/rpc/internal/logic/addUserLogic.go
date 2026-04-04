package logic

import (
	"context"

	"amigo-api/app/user/model"
	"amigo-api/app/user/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddUserLogic {
	return &AddUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddUserLogic) AddUser(in *pb.AddUserReq) (*pb.LoginSuccessResp, error) {
	// 创建数据模型
	var m model.User
	if err := copier.Copy(&m, in); err != nil {
		return nil, err
	}

	// 设置默认值
	if m.IsDelete == 0 {
		m.IsDelete = 2 // 2表示未删除
	}

	// 检查重复
	isDuplicate, err := l.svcCtx.UserModel.CheckDuplicate(l.ctx, &m)
	if err != nil {
		return nil, err
	}
	if isDuplicate {
		return nil, model.ErrDuplicate
	}

	// 插入数据
	result, err := l.svcCtx.UserModel.Insert(l.ctx, &m)
	if err != nil {
		return nil, err
	}

	// 获取插入的ID
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	m.UserId = uint64(id)

	// 构造响应
	var resp pb.LoginSuccessResp
	if err := copier.Copy(&resp, &m); err != nil {
		return nil, err
	}

	return &resp, nil
}
