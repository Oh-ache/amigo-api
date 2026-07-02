package logic

import (
	"amigo-api/app/device/model"
	"amigo-api/common/pb"
)

// toPbDeviceEvent 将 model.DeviceEvent 转换为 pb.DeviceEventResp
func toPbDeviceEvent(m *model.DeviceEvent) *pb.DeviceEventResp {
	if m == nil {
		return nil
	}
	resp := &pb.DeviceEventResp{
		DeviceEventId: m.DeviceEventId,
		DeviceId:      m.DeviceId,
		EventType:     m.EventType,
		EventLevel:    m.EventLevel,
		Title:         m.Title,
		Source:        m.Source,
		IsDelete:      m.IsDelete,
		CreateTime:    m.CreateTime,
	}
	if m.Description.Valid {
		resp.Description = m.Description.String
	}
	if m.ExtraData.Valid {
		resp.ExtraData = m.ExtraData.String
	}
	return resp
}
