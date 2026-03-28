package courier

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"pkg"

	"amigo-api/common/utils"

	"github.com/valyala/fasthttp"
)

type KuaiDiNiaoConfig struct {
	ApiKey      string
	EBusinessID string
}

type KuaiDiBiaoBody struct {
	RequestData string `json:"RequestData"`
	EBusinessID string `json:"EBusinessID"`
	RequestType string `json:"RequestType"`
	DataType    string `json:"DataType,default=2"`
	DataSign    string `json:"DataSign"`
}

type SearchReq struct {
	LogisticCode string `json:"LogisticCode"`
	CustomerName string `json:"CustomerName,optional"`
	ShipperCode  string `json:"ShipperCode"`
}

type SearchResp struct{}

type SearchResult struct {
	EBusinessID  string `json:"EBusinessID"`
	ShipperCode  string `json:"ShipperCode"`
	Success      bool   `json:"Success"`
	LogisticCode string `json:"LogisticCode"`
	State        string `json:"State"`
	StateEx      string `json:"StateEx"`
	Location     string `json:"Location"`
	Traces       []struct {
		AcceptTime    string `json:"AcceptTime"`
		AcceptStation string `json:"AcceptStation"`
		Location      string `json:"Location"`
		Action        string `json:"Action"`
	} `json:"Traces"`
}

func getSign(data, apiKey string) string {
	str := fmt.Sprintf("%s%s", data, apiKey)

	md5Str := utils.Md5(str)

	base64Str := base64.StdEncoding.EncodeToString(pkg.String2Bytes(md5Str))

	return url.QueryEscape(base64Str)
}

func Get(config *KuaiDiNiaoConfig, params *SearchReq) (*SearchResult, error) {
	url := "https://api.kdniao.com/api/dist"

	params.ShipperCode = "SF"
	byteStr, _ := json.Marshal(params)
	// body := KuaiDiBiaoBody{}
	// body.RequestData = common.Bytes2String(byteStr)
	// body.EBusinessID = config.EBusinessID
	// body.RequestType = "8001"
	// body.DataType = "2"
	// body.DataSign = getSign(body.RequestData, config.ApiKey)
	body := fmt.Sprintf("%s=%s&%s=%s&%s=%s&%s=%s&%s=%s",
		"RequestData", utils.Bytes2String(byteStr),
		"EBusinessID", config.EBusinessID,
		"RequestType", "8001",
		"DataType", "2",
		"DataSign", getSign(pkg.Bytes2String(byteStr), config.ApiKey))

	result := &SearchResult{}

	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	if err := utils.FastWithUrlencodeDo(result, fasthttp.MethodPost, url, nil, body, headers); err != nil {
		return nil, err
	}

	return result, nil
}
