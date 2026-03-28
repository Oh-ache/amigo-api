package ip

import (
	"fmt"
	"time"

	"amigo-api/common/utils"

	"github.com/valyala/fasthttp"
)

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

type Ip2AddressReq struct {
	Ip          string
	AppId       string
	AppSecurity string
}

type Ip2AddressResp struct {
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Data    struct {
		OrderNo string `json:"orderNo"`
		Result  struct {
			Continent string `json:"continent"`
			Owner     string `json:"owner"`
			Country   string `json:"country"`
			Lng       string `json:"lng"`
			Adcode    string `json:"adcode"`
			City      string `json:"city"`
			Timezone  string `json:"timezone"`
			Isp       string `json:"isp"`
			Accuracy  string `json:"accuracy"`
			Source    string `json:"source"`
			Asnumber  string `json:"asnumber"`
			Areacode  string `json:"areacode"`
			Zipcode   string `json:"zipcode"`
			District  string `json:"district"`
			Radius    string `json:"radius"`
			Prov      string `json:"prov"`
			Lat       string `json:"lat"`
		} `json:"result"`
	} `json:"data"`
}

func Ip2Address(req *Ip2AddressReq) (*IpAddress, error) {
	uri := "https://api.shumaidata.com/v4/ip/district/query"

	timestamp := time.Now().UnixMilli()

	md5Str := fmt.Sprintf("%s&%d&%s", req.AppId, timestamp, req.AppSecurity)
	sign := utils.Md5(md5Str)

	params := map[string]string{}

	params["ip"] = req.Ip
	params["appid"] = req.AppId
	params["timestamp"] = fmt.Sprintf("%d", timestamp)
	params["sign"] = sign

	result := &Ip2AddressResp{}
	if err := utils.FastWithDo(result, fasthttp.MethodGet, uri, params, nil, nil); err != nil {
		return nil, err
	}

	if result.Success != true {
		return nil, fmt.Errorf(result.Msg)
	}

	resp := &IpAddress{}
	resp.Ip = req.Ip
	resp.Country = result.Data.Result.Country
	resp.Province = result.Data.Result.Prov
	resp.City = result.Data.Result.City
	resp.Area = result.Data.Result.District
	resp.Isp = result.Data.Result.Isp
	resp.Lng = result.Data.Result.Lng
	resp.Lat = result.Data.Result.Lat

	return resp, nil
}
