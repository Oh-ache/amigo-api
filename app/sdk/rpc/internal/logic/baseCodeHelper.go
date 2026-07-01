package logic

import (
	"context"

	"amigo-api/app/baseCode/rpc/basecode"
	"amigo-api/common/pb"
)

func GetBaseCode(ctx context.Context, baseCode basecode.BaseCode, sortKey, key string) string {
	item, _ := baseCode.GetBaseCode(ctx, &pb.GetBaseCodeReq{
		SortKey: sortKey,
		Key:     key,
	})
	if item != nil {
		return item.Content
	}
	return ""
}