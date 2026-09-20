package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"amigo-api/app/sdk/rpc/internal/svc"
	"amigo-api/common/pb"
	"amigo-api/common/utils"
	"amigo-api/common/utils/plug/weather"

	"github.com/zeromicro/go-zero/core/logx"
)

// 天气缓存前缀与每日缓存 TTL
const (
	weatherCachePrefix = "cache:amigo:sdk:weather:"
	areaCodeSortKey    = "sdk"
	areaCodeItemKey    = "gaode.cityCode"
)

type GetWeatherLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetWeatherLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWeatherLogic {
	return &GetWeatherLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetWeatherLogic) GetWeather(in *pb.GetWeatherReq) (*pb.GetWeatherResp, error) {
	resp := &pb.GetWeatherResp{}

	// 1. 解析 adcode：code 与 name 二选一
	adcode := strings.TrimSpace(in.Code)
	if adcode == "" && strings.TrimSpace(in.Name) != "" {
		resolved, err := l.resolveAdcodeByName(strings.TrimSpace(in.Name))
		if err != nil {
			return nil, err
		}
		adcode = resolved
	}
	if adcode == "" {
		return nil, fmt.Errorf("code 和 name 不能同时为空")
	}

	// 2. 获取高德 key
	key, err := l.getGaodeWeatherKey()
	if err != nil {
		return nil, err
	}

	// 3. 构造 fallback 链：区 → 市 → 省
	chain := buildFallbackChain(adcode)
	extensions := "base"
	if in.Type == "2" {
		extensions = "all"
	}

	// 4. 沿链查询（带 Redis 缓存）
	var weatherResp *weather.GaodeData
	for _, code := range chain {
		data, queryErr := l.queryWeatherWithCache(key, code, extensions)
		if queryErr != nil {
			l.Logger.Errorf("查询天气失败 adcode=%s: %v", code, queryErr)
			continue
		}
		if data == nil {
			continue
		}
		if extensions == "base" && len(data.Lives) > 0 {
			weatherResp = data
			break
		}
		if extensions == "all" && len(data.Forecasts) > 0 {
			weatherResp = data
			break
		}
	}

	// 5. 组装响应（与旧逻辑保持兼容）
	if extensions == "base" && weatherResp != nil && len(weatherResp.Lives) > 0 {
		live := weatherResp.Lives[0]
		resp.Date = live.Reporttime
		resp.Week = utils.WeekdayInChinese(time.Now().Weekday())
		resp.Weather = live.Weather
		resp.Temp = live.Temperature
		resp.Wind = fmt.Sprintf("%s风%s", live.Winddirection, live.Windpower)
		resp.Wind = strings.Replace(resp.Wind, "≤", "1-", -1)
		resp.Humidity = live.Humidity
		resp.Items = make([]*pb.GetWeatherItem, 0)
	} else {
		resp.Items = make([]*pb.GetWeatherItem, 0)
	}

	return resp, nil
}

// resolveAdcodeByName 通过中文名从 baseCodeItem 表查找 adcode。
// 重名时取第一条；未找到则返回错误。
func (l *GetWeatherLogic) resolveAdcodeByName(name string) (string, error) {
	listResp, err := l.svcCtx.BaseCodeRpc.ListBaseCodeItem(l.ctx, &pb.ListBaseCodeItemReq{
		SortKey:   areaCodeSortKey,
		Key:       areaCodeItemKey,
		Content1:  name,
		IsDelete:  2,
		Page:      1,
		PageSize:  1,
	})
	if err != nil {
		return "", fmt.Errorf("查询区域字典失败: %w", err)
	}
	if listResp == nil || len(listResp.List) == 0 {
		return "", fmt.Errorf("未找到地区: %s", name)
	}
	return strconv.FormatUint(listResp.List[0].ForeignId, 10), nil
}

// getGaodeWeatherKey 从 baseCode 取高德天气 key。
func (l *GetWeatherLogic) getGaodeWeatherKey() (string, error) {
	item, err := l.svcCtx.BaseCodeRpc.GetBaseCode(l.ctx, &pb.GetBaseCodeReq{
		SortKey: areaCodeSortKey,
		Key:     "gaode.weather.key",
	})
	if err != nil {
		return "", fmt.Errorf("获取高德天气 key 失败: %w", err)
	}
	if item == nil {
		return "", fmt.Errorf("未配置 gaode.weather.key")
	}
	return item.Content, nil
}

// buildFallbackChain 构造 fallback 链：区 → 市 → 省。
// 入参要求为 6 位字符串 adcode；末尾 00 表示市/省，不再向父级收敛。
func buildFallbackChain(adcode string) []string {
	chain := make([]string, 0, 3)
	if adcode == "" {
		return chain
	}
	chain = append(chain, adcode)
	if len(adcode) != 6 {
		return chain
	}
	// 区/县级 → 市级
	if adcode[4:6] != "00" {
		city := adcode[:4] + "00"
		if city != adcode {
			chain = append(chain, city)
		}
		// 市级 → 省级（仅对 4 个直辖市/特别行政区有意义；其它省市已是省级）
		province := adcode[:2] + "0000"
		if province != adcode && province != city {
			chain = append(chain, province)
		}
	}
	return chain
}

// queryWeatherWithCache 带 Redis 缓存的天气查询，TTL 至当日 23:59:59。
func (l *GetWeatherLogic) queryWeatherWithCache(key, adcode, extensions string) (*weather.GaodeData, error) {
	cacheKey := weatherCachePrefix + adcode + ":" + extensions

	// 读缓存
	if cached, _ := l.svcCtx.RedisClient.Get(cacheKey); cached != "" {
		var data weather.GaodeData
		if err := json.Unmarshal([]byte(cached), &data); err == nil {
			return &data, nil
		}
	}

	// 调高德
	data, err := weather.GetWeather(key, adcode, extensions)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	// 写缓存（仅 TTL > 0 时写入，避免边界场景下 setex 报错）
	if ttl := secondsUntilEndOfDay(); ttl > 0 {
		if payload, marshalErr := json.Marshal(data); marshalErr == nil {
			l.svcCtx.RedisClient.Set(cacheKey, string(payload))
			l.svcCtx.RedisClient.Expire(cacheKey, ttl)
		}
	}
	return data, nil
}

// secondsUntilEndOfDay 返回距当日 23:59:59 的秒数。
func secondsUntilEndOfDay() int {
	now := time.Now()
	end := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
	secs := int(end.Sub(now).Seconds())
	if secs < 1 {
		return 0
	}
	return secs
}