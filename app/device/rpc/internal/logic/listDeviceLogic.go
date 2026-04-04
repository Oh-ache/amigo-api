package logic

import (
	"context"

	"amigo-api/app/device/model"
	"amigo-api/app/device/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListDeviceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListDeviceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDeviceLogic {
	return &ListDeviceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListDeviceLogic) ListDevice(in *pb.ListDeviceReq) (*pb.ListDeviceResp, error) {
	search := &model.DeviceSearch{}
	copier.Copy(search, in)

	list, total, err := l.svcCtx.DeviceModel.List(l.ctx, search)
	if err != nil {
		return nil, err
	}

	var respList []*pb.DeviceResp
	for _, item := range list {
		var respItem pb.DeviceResp
		if err := copier.Copy(&respItem, item); err != nil {
			return nil, err
		}
		respList = append(respList, &respItem)
	}

	return &pb.ListDeviceResp{
		List:  respList,
		Total: total,
	}, nil
}
