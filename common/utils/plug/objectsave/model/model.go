package model

type UploadResult struct {
	FileUrl    string `json:"file_url"`
	FileSize   int64  `json:"file_size"`
	UploadTime int64  `json:"upload_time"`
}

type UploadToken struct {
	Token      string `json:"token"`
	ExpireTime int64  `json:"expire_time"`
}

type OssConfig struct {
	Endpoint        string `json:"endpoint"`
	AccessKeyId     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	Bucket          string `json:"bucket"`
	Region          string `json:"region"`
}

type StorageConfig struct {
	Type      string
	OssConfig *OssConfig
	Extra     map[string]string
}

type CosConfig struct {
	Region    string `json:"region"`
	Bucket    string `json:"bucket"`
	SecretId  string `json:"secret_id"`
	SecretKey string `json:"secret_key"`
	AppId     string `json:"app_id"`
}
