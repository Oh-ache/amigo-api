package chat

import (
	"fmt"

	"amigo-api/common/utils"

	"github.com/valyala/fasthttp"
)

type XingHuoResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Sid     string `json:"sid"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		Index int `json:"index"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

type XinghuoReq struct {
	Model    string            `json:"model"`
	Messages []XinghuoMessages `json:"messages"`
	Stream   bool              `json:"stream"`
}
type XinghuoMessages struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func XinghuoChat(passwd, model, content string) (string, error) {
	requestUrl := "https://spark-api-open.xf-yun.com/v1/chat/completions"
	requestPasswd := passwd

	body := &XinghuoReq{
		Model: model,
		Messages: []XinghuoMessages{
			{
				Role:    "user",
				Content: content,
			},
		},
		Stream: false,
	}

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = fmt.Sprintf("Bearer %s", requestPasswd)

	result := &XingHuoResp{}
	if err := utils.FastWithDo(result, fasthttp.MethodPost, requestUrl, nil, body, headers); err != nil {
		return "", err
	}
	return result.Choices[0].Message.Content, nil
}
