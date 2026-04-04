package logic

import (
	"context"

	"amigo-api/app/baseCode/model"

	"amigo-api/app/baseCode/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddBaseCodeSortLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddBaseCodeSortLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddBaseCodeSortLogic {
	return &AddBaseCodeSortLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddBaseCodeSortLogic) AddBaseCodeSort(in *pb.AddBaseCodeSortReq) (*pb.BaseCodeSortResp, error) {
	// 创建 model.BaseCodeSort 实例
	var sort model.BaseCodeSort

	// 使用 copier 库将入参复制到 sort 实例
	if err := copier.Copy(&sort, in); err != nil {
		l.Errorf("Failed to copy AddBaseCodeSortReq to BaseCodeSort: %v", err)
		return nil, err
	}

	// 设置默认值
	sort.IsDelete = 2 // 2表示未删除

	// 先判断是否重复
	isDuplicate, err := l.svcCtx.BaseCodeSortModel.CheckDuplicate(l.ctx, &sort)
	if err != nil {
		l.Errorf("Failed to check duplicate for BaseCodeSort: %v", err)
		return nil, err
	}
	if isDuplicate {
		return nil, model.ErrDuplicate // 直接使用 model 中定义的 ErrDuplicate
	}

	// 写入数据库
	result, err := l.svcCtx.BaseCodeSortModel.Insert(l.ctx, &sort)
	if err != nil {
		l.Errorf("Failed to insert BaseCodeSort: %v", err)
		return nil, err
	}

	// 获取插入后的主键ID
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		l.Errorf("Failed to get last insert ID: %v", err)
		return nil, err
	}
	sort.BaseCodeSortId = uint64(lastInsertId)

	// 创建返回响应
	var resp pb.BaseCodeSortResp
	if err := copier.Copy(&resp, &sort); err != nil {
		l.Errorf("Failed to copy BaseCodeSort to BaseCodeSortResp: %v", err)
		return nil, err
	}

	return &resp, nil
}
