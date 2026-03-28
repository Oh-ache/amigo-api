package logic

import (
	"context"
	"fmt"
	"strings"
	"time"

	"amigo-api/app/sdk/rpc/internal/svc"
	"amigo-api/common/pb"
	"amigo-api/common/utils"
	"amigo-api/common/utils/plug/weather"

	"github.com/zeromicro/go-zero/core/logx"
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
	// TODO 加缓存减少请求次数
	resp := &pb.GetWeatherResp{}

	getBaseCodeReq := &pb.GetBaseCodeReq{}
	getBaseCodeReq.SortKey = "sdk"
	item := &pb.BaseCodeResp{}

	getBaseCodeReq.Key = "gaode.weather.key"
	item, _ = l.svcCtx.BaseCodeRpc.GetBaseCode(l.ctx, getBaseCodeReq)
	key := item.Content

	extensions := "base"
	if in.Type == "2" {
		extensions = "all"
	}

	weatherResp, _ := weather.GetWeather(key, in.Code, extensions)
	if extensions == "base" {
		resp.Date = weatherResp.Lives[0].Reporttime
		resp.Week = utils.WeekdayInChinese(time.Now().Weekday())
		resp.Weather = weatherResp.Lives[0].Weather
		resp.Temp = weatherResp.Lives[0].Temperature
		resp.Wind = fmt.Sprintf("%s风%s", weatherResp.Lives[0].Winddirection, weatherResp.Lives[0].Windpower)
		resp.Wind = strings.Replace(resp.Wind, "≤", "1-", -1)
		resp.Humidity = weatherResp.Lives[0].Humidity
		resp.Items = make([]*pb.GetWeatherItem, 0)
	}

	return resp, nil
}
