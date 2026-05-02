package mqueue

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"amigo-api/common/mqueue"
	"amigo-api/common/pb"
	"amigo-api/common/utils"
	"amigo-api/common/utils/plug/message"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
)

var RedisClient *redis.Client

func InitRedis(client *redis.Client) {
	RedisClient = client
}

type BaseCodeRpcClient interface {
	GetBaseCode(ctx context.Context, req *pb.GetBaseCodeReq) (*pb.BaseCodeResp, error)
}

type AiRpcClient interface {
	UpdateTask(ctx context.Context, req *pb.UpdateTaskReq) (*pb.UpdateTaskResp, error)
	UploadUrl(ctx context.Context, req *pb.UploadUrlReq) (*pb.UploadUrlResp, error)
}

type SendSmsHandler struct{}

func NewSendSmsHandler() *SendSmsHandler {
	return &SendSmsHandler{}
}

func (h *SendSmsHandler) Name() string {
	return "send_sms"
}

func (h *SendSmsHandler) Handle(ctx context.Context, task *mqueue.Task) error {
	dataMap, ok := task.Data["data"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid data format")
	}

	dataBytes, err := json.Marshal(dataMap)
	if err != nil {
		return err
	}

	var data message.PushContext
	if err := json.Unmarshal(dataBytes, &data); err != nil {
		return err
	}

	sendType, _ := task.Data["send_type"].(string)
	code, _ := task.Data["code"].(string)

	if err := message.PushMessage(&data); err != nil {
		return err
	}

	redisKey := fmt.Sprintf("%s%s:%s", utils.SEND_CODE_KEY, sendType, data.Mobile)
	RedisClient.Set(ctx, redisKey, code, 180*time.Second)

	logx.Infof("SendSmsHandler: SMS sent to %s, code: %s", data.Mobile, code)
	return nil
}

type AiTaskHandler struct {
	baseCodeRpc BaseCodeRpcClient
	aiRpcClient AiRpcClient
	httpClient  *http.Client
}

func NewAiTaskHandler() *AiTaskHandler {
	return &AiTaskHandler{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (h *AiTaskHandler) Name() string {
	return "ai_task"
}

func (h *AiTaskHandler) SetBaseCodeRpc(rpc BaseCodeRpcClient) {
	h.baseCodeRpc = rpc
}

func (h *AiTaskHandler) SetAiRpcClient(rpc AiRpcClient) {
	h.aiRpcClient = rpc
}

func (h *AiTaskHandler) Handle(ctx context.Context, task *mqueue.Task) error {
	id, ok := task.Data["id"].(float64)
	if !ok {
		return fmt.Errorf("invalid id")
	}

	taskType, _ := task.Data["task_type"].(string)
	prompt, _ := task.Data["prompt"].(string)
	params, _ := task.Data["params"].(string)

	logx.Infof("AiTaskHandler: processing task id=%.0f, type=%s", id, taskType)

	switch taskType {
	case "text_to_image":
		return h.handleTextToImage(ctx, int64(id), prompt, params)
	case "video":
		return h.handleVideo(ctx, int64(id), prompt, params)
	case "audio":
		return h.handleAudio(ctx, int64(id), prompt, params)
	default:
		return fmt.Errorf("unknown task type: %s", taskType)
	}
}

func (h *AiTaskHandler) handleTextToImage(ctx context.Context, id int64, prompt, params string) error {
	logx.Infof("AiTaskHandler: text_to_image task id=%d, calling MiniMax API", id)

	token, err := h.getMiniMaxToken(ctx)
	if err != nil {
		logx.Errorf("AiTaskHandler: failed to get token: %v", err)
		return h.updateTaskError(ctx, id, err.Error())
	}

	reqBody := map[string]interface{}{
		"model":           "image-01",
		"prompt":          prompt,
		"response_format": "url",
	}
	reqBodyBytes, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", "https://api.minimaxi.com/v1/image_generation", bytes.NewReader(reqBodyBytes))
	if err != nil {
		return h.updateTaskError(ctx, id, err.Error())
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		logx.Errorf("AiTaskHandler: MiniMax API call failed: %v", err)
		return h.updateTaskError(ctx, id, err.Error())
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return h.updateTaskError(ctx, id, err.Error())
	}

	responseInfo, _ := json.Marshal(result)

	if baseResp, ok := result["base_resp"].(map[string]interface{}); ok {
		if statusCode, ok := baseResp["status_code"].(float64); ok && statusCode != 0 {
			msg, _ := baseResp["status_msg"].(string)
			h.updateTaskResponseInfo(ctx, id, string(responseInfo))
			return h.updateTaskError(ctx, id, msg)
		}
	}

	var imageUrls []string
	if data, ok := result["data"].(map[string]interface{}); ok {
		if urls, ok := data["image_urls"].([]interface{}); ok {
			for _, u := range urls {
				if urlStr, ok := u.(string); ok {
					imageUrls = append(imageUrls, urlStr)
				}
			}
		}
	}

	if len(imageUrls) > 0 {
		ossUrl, err := h.uploadToOss(ctx, imageUrls[0], "text_to_image")
		if err != nil {
			logx.Errorf("AiTaskHandler: upload to oss failed: %v", err)
			h.updateTaskResponseInfo(ctx, id, string(responseInfo))
			return h.updateTaskError(ctx, id, err.Error())
		}

		h.updateTaskSuccess(ctx, id, string(responseInfo), ossUrl)
		logx.Infof("AiTaskHandler: text_to_image task id=%d completed, result_url=%s", id, ossUrl)
	}

	return nil
}

func (h *AiTaskHandler) handleVideo(ctx context.Context, id int64, prompt, params string) error {
	logx.Infof("AiTaskHandler: video task id=%d, creating MiniMax video task", id)

	token, err := h.getMiniMaxToken(ctx)
	if err != nil {
		return h.updateTaskError(ctx, id, err.Error())
	}

	createResp, err := h.createVideoTask(ctx, token, prompt, params)
	if err != nil {
		return h.updateTaskError(ctx, id, err.Error())
	}

	taskId := createResp.TaskId
	logx.Infof("AiTaskHandler: video task created, task_id=%s, id=%d", taskId, id)

	h.updateTaskTaskId(ctx, id, taskId, createResp.ResponseInfo)

	time.Sleep(5 * time.Second)

	maxAttempts := 150
	for i := 0; i < maxAttempts; i++ {
		status, fileId, err := h.queryVideoTask(ctx, token, taskId)
		if err != nil {
			logx.Errorf("AiTaskHandler: query video task failed: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		if status == "success" && fileId != "" {
			downloadUrl, err := h.downloadVideo(ctx, token, fileId)
			if err != nil {
				return h.updateTaskError(ctx, id, err.Error())
			}

			ossUrl, err := h.uploadToOss(ctx, downloadUrl, "video")
			if err != nil {
				return h.updateTaskError(ctx, id, err.Error())
			}

			h.updateTaskSuccess(ctx, id, createResp.ResponseInfo, ossUrl)
			logx.Infof("AiTaskHandler: video task id=%d completed", id)
			return nil
		}

		if status == "failed" {
			return h.updateTaskError(ctx, id, "video generation failed")
		}

		time.Sleep(2 * time.Second)
	}

	return h.updateTaskError(ctx, id, "video generation timeout")
}

func (h *AiTaskHandler) handleAudio(ctx context.Context, id int64, prompt, params string) error {
	logx.Infof("AiTaskHandler: audio task id=%d, creating MiniMax audio task", id)

	token, err := h.getMiniMaxToken(ctx)
	if err != nil {
		return h.updateTaskError(ctx, id, err.Error())
	}

	createResp, err := h.createAudioTask(ctx, token, prompt, params)
	if err != nil {
		return h.updateTaskError(ctx, id, err.Error())
	}

	taskId := createResp.TaskId
	logx.Infof("AiTaskHandler: audio task created, task_id=%s, id=%d", taskId, id)

	h.updateTaskTaskId(ctx, id, taskId, createResp.ResponseInfo)

	time.Sleep(5 * time.Second)

	maxAttempts := 150
	for i := 0; i < maxAttempts; i++ {
		status, fileId, err := h.queryAudioTask(ctx, token, taskId)
		if err != nil {
			logx.Errorf("AiTaskHandler: query audio task failed: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		if status == "Success" && fileId != "" {
			downloadUrl, err := h.downloadAudio(ctx, token, fileId)
			if err != nil {
				return h.updateTaskError(ctx, id, err.Error())
			}

			ossUrl, err := h.uploadToOss(ctx, downloadUrl, "audio")
			if err != nil {
				return h.updateTaskError(ctx, id, err.Error())
			}

			h.updateTaskSuccess(ctx, id, createResp.ResponseInfo, ossUrl)
			logx.Infof("AiTaskHandler: audio task id=%d completed", id)
			return nil
		}

		if status == "Failed" {
			return h.updateTaskError(ctx, id, "audio generation failed")
		}

		time.Sleep(2 * time.Second)
	}

	return h.updateTaskError(ctx, id, "audio generation timeout")
}

func (h *AiTaskHandler) getMiniMaxToken(ctx context.Context) (string, error) {
	if h.baseCodeRpc == nil {
		return "", fmt.Errorf("baseCodeRpc not initialized")
	}

	req := &pb.GetBaseCodeReq{
		SortKey: "sdk",
		Key:     "minimax.tokenplankey",
	}
	resp, err := h.baseCodeRpc.GetBaseCode(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.Content, nil
}

type CreateTaskResponse struct {
	TaskId       string
	ResponseInfo string
}

func (h *AiTaskHandler) createVideoTask(ctx context.Context, token, prompt, params string) (*CreateTaskResponse, error) {
	reqBody := map[string]interface{}{
		"model":  "MiniMax-Hailuo-2.3",
		"prompt": prompt,
	}
	if params != "" {
		var p map[string]interface{}
		if err := json.Unmarshal([]byte(params), &p); err == nil {
			if duration, ok := p["duration"].(float64); ok {
				reqBody["duration"] = int(duration)
			}
			if resolution, ok := p["resolution"].(string); ok {
				reqBody["resolution"] = resolution
			}
		}
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", "https://api.minimaxi.com/v1/video_generation", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	responseInfo, _ := json.Marshal(result)

	if baseResp, ok := result["base_resp"].(map[string]interface{}); ok {
		if statusCode, ok := baseResp["status_code"].(float64); ok && statusCode != 0 {
			msg, _ := baseResp["status_msg"].(string)
			return nil, fmt.Errorf("%s", msg)
		}
	}

	taskId, _ := result["task_id"].(string)
	return &CreateTaskResponse{
		TaskId:       taskId,
		ResponseInfo: string(responseInfo),
	}, nil
}

func (h *AiTaskHandler) queryVideoTask(ctx context.Context, token, taskId string) (status string, fileId string, err error) {
	req, err := http.NewRequest("GET", "https://api.minimaxi.com/v1/video_generation/query?task_id="+taskId, nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := h.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", err
	}

	status, _ = result["status"].(string)
	fileId, _ = result["file_id"].(string)
	return status, fileId, nil
}

func (h *AiTaskHandler) downloadVideo(ctx context.Context, token, fileId string) (string, error) {
	req, err := http.NewRequest("GET", "https://api.minimaxi.com/v1/files/retrieve?file_id="+fileId, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := h.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if file, ok := result["file"].(map[string]interface{}); ok {
		if downloadUrl, ok := file["download_url"].(string); ok {
			return downloadUrl, nil
		}
	}

	return "", fmt.Errorf("download_url not found")
}

func (h *AiTaskHandler) createAudioTask(ctx context.Context, token, prompt, params string) (*CreateTaskResponse, error) {
	reqBody := map[string]interface{}{
		"model": "speech-2.8-hd",
		"text":  prompt,
		"voice_setting": map[string]interface{}{
			"voice_id": "moss_audio_ce44fc67-7ce3-11f0-8de5-96e35d26fb85",
			"speed":    1,
		},
	}
	if params != "" {
		var p map[string]interface{}
		if err := json.Unmarshal([]byte(params), &p); err == nil {
			if voiceId, ok := p["voice_id"].(string); ok {
				reqBody["voice_setting"] = map[string]interface{}{
					"voice_id": voiceId,
					"speed":    1,
				}
			}
		}
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", "https://api.minimaxi.com/v1/t2a_async_v2", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	responseInfo, _ := json.Marshal(result)

	if baseResp, ok := result["base_resp"].(map[string]interface{}); ok {
		if statusCode, ok := baseResp["status_code"].(float64); ok && statusCode != 0 {
			msg, _ := baseResp["status_msg"].(string)
			return nil, fmt.Errorf("%s", msg)
		}
	}

	taskId := ""
	if tid, ok := result["task_id"].(float64); ok {
		taskId = fmt.Sprintf("%.0f", tid)
	}
	return &CreateTaskResponse{
		TaskId:       taskId,
		ResponseInfo: string(responseInfo),
	}, nil
}

func (h *AiTaskHandler) queryAudioTask(ctx context.Context, token string, taskId string) (status string, fileId string, err error) {
	req, err := http.NewRequest("GET", "https://api.minimaxi.com/v1/query/t2a_async_query_v2?task_id="+taskId, nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := h.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", err
	}

	status, _ = result["status"].(string)
	if fid, ok := result["file_id"].(float64); ok {
		fileId = fmt.Sprintf("%.0f", fid)
	}
	return status, fileId, nil
}

func (h *AiTaskHandler) downloadAudio(ctx context.Context, token, fileId string) (string, error) {
	req, err := http.NewRequest("GET", "https://api.minimaxi.com/v1/files/retrieve?file_id="+fileId, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := h.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if file, ok := result["file"].(map[string]interface{}); ok {
		if downloadUrl, ok := file["download_url"].(string); ok {
			return downloadUrl, nil
		}
	}

	return "", fmt.Errorf("download_url not found")
}

func (h *AiTaskHandler) uploadToOss(ctx context.Context, url string, taskType string) (string, error) {
	if h.aiRpcClient == nil {
		return "", fmt.Errorf("aiRpcClient not initialized")
	}

	fileName := generateFileName(taskType)
	resp, err := h.aiRpcClient.UploadUrl(ctx, &pb.UploadUrlReq{
		FileName: fileName,
		Url:      url,
	})
	if err != nil {
		return "", err
	}

	return resp.Url, nil
}

func generateFileName(taskType string) string {
	now := time.Now()
	dateStr := now.Format("20060102150405")
	randNum := rand.Intn(10000)

	var ext string
	switch taskType {
	case "text_to_image":
		ext = ".jpg"
	case "video":
		ext = ".mp4"
	case "audio":
		ext = ".mp3"
	default:
		ext = ".bin"
	}

	return fmt.Sprintf("test/%s%d%s", dateStr, randNum, ext)
}

func (h *AiTaskHandler) updateTaskError(ctx context.Context, id int64, errorMsg string) error {
	if h.aiRpcClient == nil {
		logx.Errorf("aiRpcClient not initialized, cannot update task error")
		return fmt.Errorf(errorMsg)
	}

	_, err := h.aiRpcClient.UpdateTask(ctx, &pb.UpdateTaskReq{
		Id:       id,
		Status:   2,
		ErrorMsg: errorMsg,
	})
	if err != nil {
		logx.Errorf("updateTaskError failed: %v", err)
	}
	return fmt.Errorf(errorMsg)
}

func (h *AiTaskHandler) updateTaskSuccess(ctx context.Context, id int64, responseInfo, resultUrl string) {
	if h.aiRpcClient == nil {
		logx.Errorf("aiRpcClient not initialized, cannot update task success")
		return
	}

	_, err := h.aiRpcClient.UpdateTask(ctx, &pb.UpdateTaskReq{
		Id:           id,
		Status:       1,
		ResponseInfo: responseInfo,
		ResultUrl:    resultUrl,
	})
	if err != nil {
		logx.Errorf("updateTaskSuccess failed: %v", err)
	}
}

func (h *AiTaskHandler) updateTaskTaskId(ctx context.Context, id int64, taskId, responseInfo string) {
	if h.aiRpcClient == nil {
		logx.Errorf("aiRpcClient not initialized, cannot update task_id")
		return
	}

	_, err := h.aiRpcClient.UpdateTask(ctx, &pb.UpdateTaskReq{
		Id:           id,
		TaskId:       taskId,
		ResponseInfo: responseInfo,
	})
	if err != nil {
		logx.Errorf("updateTaskTaskId failed: %v", err)
	}
}

func (h *AiTaskHandler) updateTaskResponseInfo(ctx context.Context, id int64, responseInfo string) {
	if h.aiRpcClient == nil {
		return
	}

	_, err := h.aiRpcClient.UpdateTask(ctx, &pb.UpdateTaskReq{
		Id:           id,
		ResponseInfo: responseInfo,
	})
	if err != nil {
		logx.Errorf("updateTaskResponseInfo failed: %v", err)
	}
}
