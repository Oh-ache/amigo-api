package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"amigo-api/app/sdk/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/valyala/fasthttp"
	"github.com/zeromicro/go-zero/core/logx"
)

type ImageGenerationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewImageGenerationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ImageGenerationLogic {
	return &ImageGenerationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// MiniMaxImageReq MiniMax 文生图请求参数
type MiniMaxImageReq struct {
	Model          string `json:"model"`
	Prompt         string `json:"prompt"`
	N              int    `json:"n"`
	Size           string `json:"size"`
	ResponseFormat string `json:"response_format"`
}

// MiniMaxImageResp MiniMax 文生图响应
type MiniMaxImageResp struct {
	BaseResp struct {
		StatusCode    int    `json:"status_code"`
		StatusMessage string `json:"status_message"`
	} `json:"base_resp"`
	Data struct {
		ImageUrls []string `json:"image_urls"`
	} `json:"data"`
}

// miniMaxRequest 独立的请求方法，设置更长的超时时间
func miniMaxRequest(result any, uri string, token string, body any) error {
	client := &fasthttp.Client{
		ReadTimeout:  time.Minute * 2,
		WriteTimeout: time.Minute * 2,
	}
	req, resp := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
	}()
	req.Header.SetMethod("POST")
	req.SetRequestURI(uri)
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	if body != nil {
		marshal, _ := json.Marshal(body)
		req.SetBody(marshal)
	}

	if err := client.Do(req, resp); err != nil {
		return err
	}

	fmt.Println(string(resp.Body()))
	return json.Unmarshal(resp.Body(), result)
}

func (l *ImageGenerationLogic) ImageGeneration(in *pb.BaseAiReq) (*pb.AiImageGenerationResp, error) {
	// 从数据库获取 minimax token
	token := GetBaseCode(l.ctx, l.svcCtx.BaseCodeRpc, "sdk", "minimax.tokenplankey")
	fmt.Println(token)

	// MiniMax API 地址
	uri := "https://api.minimax.chat/v1/image_generation"

	// 构建请求体
	reqBody := MiniMaxImageReq{
		Model:          "image-01",
		Prompt:         in.Prompt,
		N:              1,
		Size:           "1024x1024",
		ResponseFormat: "url",
	}

	var resp MiniMaxImageResp
	err := miniMaxRequest(&resp, uri, token, reqBody)
	if err != nil {
		return nil, err
	}
	fmt.Println(resp)

	if resp.BaseResp.StatusCode != 0 {
		logx.Errorf("MiniMax API error: %s", resp.BaseResp.StatusMessage)
		return nil, nil
	}

	imageUrl := ""

	if len(resp.Data.ImageUrls) > 0 {
		imageUrl = resp.Data.ImageUrls[0]
	}

	return &pb.AiImageGenerationResp{
		Url: imageUrl,
	}, nil
}
