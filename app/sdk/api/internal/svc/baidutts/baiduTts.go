package baidutts

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"amigo-api/app/baseCode/rpc/basecode"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	tokenURL = "https://aip.baidubce.com/oauth/2.0/token"
	synthURL = "https://tsn.baidu.com/text2audio"

	apiKeyCode    = "baidu..apikey"
	secretKeyCode = "baidu.secretkey"

	tokenSkewSec = 86400
	cacheTTL     = 5 * time.Minute
)

type Client struct {
	baseCodeRpc basecode.BaseCode
	http        *http.Client

	mu          sync.Mutex
	apiKey      string
	secretKey   string
	keyLoadedAt time.Time

	tokenMu   sync.Mutex
	token     string
	tokenExp  time.Time
}

// New 创建客户端。baseCodeRpc 用于从 base_code 字典读 API Key/Secret。
func New(bc basecode.BaseCode) *Client {
	return &Client{
		baseCodeRpc: bc,
		http:        &http.Client{Timeout: 30 * time.Second},
	}
}

// loadKeys 从 base_code 表读取 API Key/Secret（缓存 5 分钟）。
func (c *Client) loadKeys(ctx context.Context) (string, string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.apiKey != "" && c.secretKey != "" && time.Since(c.keyLoadedAt) < cacheTTL {
		return c.apiKey, c.secretKey, nil
	}
	apiKey, err := c.fetchKey(ctx, apiKeyCode)
	if err != nil {
		return "", "", fmt.Errorf("读 %s: %w", apiKeyCode, err)
	}
	secretKey, err := c.fetchKey(ctx, secretKeyCode)
	if err != nil {
		return "", "", fmt.Errorf("读 %s: %w", secretKeyCode, err)
	}
	c.apiKey = apiKey
	c.secretKey = secretKey
	c.keyLoadedAt = time.Now()
	return apiKey, secretKey, nil
}

func (c *Client) fetchKey(ctx context.Context, key string) (string, error) {
	resp, err := c.baseCodeRpc.GetBaseCode(ctx, &pb.GetBaseCodeReq{
		SortKey: "sdk",
		Key:     key,
	})
	if err != nil {
		return "", err
	}
	if resp.Content == "" {
		return "", fmt.Errorf("base_code[sort_key=sdk, key=%s] content 为空", key)
	}
	return resp.Content, nil
}

// getToken 获取 access_token，自动缓存 + 提前 1 天刷新。
func (c *Client) getToken(ctx context.Context, apiKey, secretKey string) (string, error) {
	c.tokenMu.Lock()
	defer c.tokenMu.Unlock()
	if c.token != "" && time.Now().Before(c.tokenExp) {
		return c.token, nil
	}
	q := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {apiKey},
		"client_secret": {secretKey},
	}
	tokenReq, err := http.NewRequestWithContext(ctx, http.MethodGet, tokenURL+"?"+q.Encode(), nil)
	if err != nil {
		return "", err
	}
	resp, err := c.http.Do(tokenReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var r struct {
		AccessToken   string `json:"access_token"`
		ExpiresIn     int    `json:"expires_in"`
		Error         string `json:"error"`
		ErrorDesc     string `json:"error_description"`
	}
	if err := json.Unmarshal(body, &r); err != nil {
		return "", fmt.Errorf("token JSON: %v (body=%s)", err, string(body))
	}
	if r.Error != "" {
		return "", errors.New(r.Error + ": " + r.ErrorDesc)
	}
	if r.AccessToken == "" {
		return "", errors.New("access_token 为空: " + string(body))
	}

	c.token = r.AccessToken
	effective := r.ExpiresIn - tokenSkewSec
	if effective < 3600 {
		effective = 3600
	}
	c.tokenExp = time.Now().Add(time.Duration(effective) * time.Second)
	logx.WithContext(ctx).Infof("百度 TTS token 刷新成功 (有效 %d s)", effective)
	return c.token, nil
}

// Synthesize 调用百度短文本在线合成，返回 PCM 字节。
func (c *Client) Synthesize(ctx context.Context, text string, per, spd, pit, vol int) ([]byte, error) {
	apiKey, secretKey, err := c.loadKeys(ctx)
	if err != nil {
		return nil, err
	}
	tok, err := c.getToken(ctx, apiKey, secretKey)
	if err != nil {
		return nil, fmt.Errorf("token: %w", err)
	}

	form := url.Values{
		"tex":  {text},
		"tok":  {tok},
		"cuid": {"amigo-server"},
		"ctp":  {"1"},
		"lan":  {"zh"},
		"aue":  {"4"},
		"spd":  {strconv.Itoa(spd)},
		"pit":  {strconv.Itoa(pit)},
		"vol":  {strconv.Itoa(vol)},
		"per":  {strconv.Itoa(per)},
	}
	resp, err := c.http.PostForm(synthURL, form)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	ct := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "audio") {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("百度 TTS 异常 CT=%s body=%s", ct, string(b))
	}
	audio, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	logx.WithContext(ctx).Infof("百度 TTS 合成成功 %d 字节 (%.1f 秒)", len(audio), float64(len(audio)/2)/16000.0)
	return audio, nil
}