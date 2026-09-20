package logic

import (
	"context"
	"encoding/json"

	"amigo-api/app/sdk/rpc/internal/svc"
	"amigo-api/common/pb"
	"amigo-api/common/utils/plug/ip"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

// IP 解析结果 Redis 缓存
const (
	ipCachePrefix = "cache:amigo:sdk:ip:"
	ipCacheTTL    = 24 * 3600 // 24 小时
)

type IpToAddressLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIpToAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IpToAddressLogic {
	return &IpToAddressLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IpToAddressLogic) IpToAddress(in *pb.IpToAddressReq) (*pb.IpToAddressResp, error) {
	cacheKey := ipCachePrefix + in.Ip

	// 1. 优先读 Redis 缓存
	if cached, _ := l.svcCtx.RedisClient.Get(cacheKey); cached != "" {
		var resp pb.IpToAddressResp
		if err := json.Unmarshal([]byte(cached), &resp); err == nil {
			return &resp, nil
		}
	}

	// 2. 缓存未命中，调用 ip9.com.cn
	param := &ip.Ip2AddressReq{Ip: in.Ip}
	result, err := ip.Ip2Address(param)
	if err != nil {
		return nil, err
	}

	resp := &pb.IpToAddressResp{}
	copier.Copy(resp, result)

	// 3. 写入缓存，TTL 24 小时
	if data, _ := json.Marshal(resp); len(data) > 0 {
		l.svcCtx.RedisClient.Set(cacheKey, string(data))
		l.svcCtx.RedisClient.Expire(cacheKey, ipCacheTTL)
	}

	return resp, nil
}