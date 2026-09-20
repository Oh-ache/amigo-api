package ip

import (
	"fmt"

	"amigo-api/common/utils"

	"github.com/valyala/fasthttp"
)

// IP 解析返回结果
type IpAddress struct {
	Ip       string
	Country  string
	Province string
	City     string
	Area     string
	Isp      string
	Lng      string
	Lat      string
}

// IP 解析入参（ip9.com.cn 接口无需鉴权）
type Ip2AddressReq struct {
	Ip string
}

// ip9.com.cn 接口原始返回结构
type Ip2AddressResp struct {
	Ret  int     `json:"ret"`
	Data struct {
		Ip          string `json:"ip"`
		Country     string `json:"country"`
		CountryCode string `json:"country_code"`
		Prov        string `json:"prov"`
		City        string `json:"city"`
		CityCode    string `json:"city_code"`
		Area        string `json:"area"`
		PostCode    string `json:"post_code"`
		AreaCode    string `json:"area_code"`
		Isp         string `json:"isp"`
		Lng         string `json:"lng"`
		Lat         string `json:"lat"`
		BigArea     string `json:"big_area"`
	} `json:"data"`
	Qt float64 `json:"qt"`
}

// Ip2Address 调用 ip9.com.cn 接口将 IP 转换为省市区信息
func Ip2Address(req *Ip2AddressReq) (*IpAddress, error) {
	uri := "https://ip9.com.cn/get"
	params := map[string]string{"ip": req.Ip}

	result := &Ip2AddressResp{}
	if err := utils.FastWithDo(result, fasthttp.MethodGet, uri, params, nil, nil); err != nil {
		return nil, err
	}
	if result.Ret != 200 {
		return nil, fmt.Errorf("ip9.com.cn 返回 ret=%d", result.Ret)
	}

	return &IpAddress{
		Ip:       req.Ip,
		Country:  result.Data.Country,
		Province: result.Data.Prov,
		City:     result.Data.City,
		Area:     result.Data.Area,
		Isp:      result.Data.Isp,
		Lng:      result.Data.Lng,
		Lat:      result.Data.Lat,
	}, nil
}