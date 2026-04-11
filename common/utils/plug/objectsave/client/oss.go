package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"amigo-api/common/utils"
	"amigo-api/common/utils/plug/objectsave/model"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

type OssClient struct {
	Client     *oss.Client
	Region     string // 存储区域
	BucketName string // 存储空间名称
	ObjectName string // 对象名称
}

func NewOssClient(ossConfig *model.OssConfig) (StorageClient, error) {
	os.Setenv("OSS_ACCESS_KEY_ID", ossConfig.AccessKeyId)
	os.Setenv("OSS_ACCESS_KEY_SECRET", ossConfig.AccessKeySecret)

	ossClient := &OssClient{
		BucketName: ossConfig.Bucket,
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithEndpoint(ossConfig.Endpoint).
		WithRegion(ossConfig.Region)
	ossClient.Client = oss.NewClient(cfg)

	return ossClient, nil
}

func (o *OssClient) UploadFile(fileName string, data []byte) (*model.UploadResult, error) {
	body := bytes.NewReader(data)

	// 创建上传对象的请求
	request := &oss.PutObjectRequest{
		Bucket: &o.BucketName, // 存储空间名称
		Key:    &fileName,     // 对象名称
		Body:   body,          // 要上传的字符串内容
	}

	// 发送上传对象的请求
	result, err := o.Client.PutObject(context.TODO(), request)
	if err != nil {
		return nil, fmt.Errorf("failed to put object: %w", err)
	}
	return &model.UploadResult{
		FileUrl:    *result.ContentMD5,
		FileSize:   0,
		UploadTime: time.Now().Unix(),
	}, nil
}

func (o *OssClient) UploadUrl(fileName string, url string) (*model.UploadResult, error) {
	// SSRF 防护：校验 URL 协议和目标 IP
	if err := utils.ValidateURL(url); err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch URL: HTTP %d", resp.StatusCode)
	}

	// 创建上传对象的请求
	request := &oss.PutObjectRequest{
		Bucket: &o.BucketName,        // 存储空间名称
		Key:    &fileName,            // 对象名称
		Body:   io.Reader(resp.Body), // 要上传的网络流内容
	}

	// 发送上传对象的请求
	result, err := o.Client.PutObject(context.TODO(), request)
	if err != nil {
		return nil, fmt.Errorf("failed to put object: %w", err)
	}
	return &model.UploadResult{
		FileUrl:    *result.ContentMD5,
		FileSize:   0,
		UploadTime: time.Now().Unix(),
	}, nil
}

func (o *OssClient) GetUploadToken(fileName string) (*model.UploadToken, error) {
	// 生成PutObject的预签名URL
	result, err := o.Client.Presign(context.TODO(), &oss.PutObjectRequest{
		Bucket:      &o.BucketName,
		Key:         &fileName,
		ContentType: oss.Ptr("image/png"), // 请确保在服务端生成该签名URL时设置的ContentType与在使用URL时设置的ContentType一致
		// StorageClass: oss.StorageClassStandard, // 请确保在服务端生成该签名URL时设置的StorageClass与在使用URL时设置的StorageClass一致
		// Metadata:    map[string]string{"key1": "value1", "key2": "value2"}, // 请确保在服务端生成该签名URL时设置的Metadata与在使用URL时设置的Metadata一致
	},
		oss.PresignExpires(900*time.Minute),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate upload token: %w", err)
	}

	return &model.UploadToken{
		Token:      result.URL,
		ExpireTime: time.Now().Add(10 * time.Minute).Unix(),
	}, nil
}

func (o *OssClient) GetStorageType() string {
	return "oss"
}
