package factory

import (
	"fmt"

	"amigo-api/common/utils/plug/objectsave/client"
	"amigo-api/common/utils/plug/objectsave/model"
)

type StorageFactory struct{}

func NewStorageFactory() *StorageFactory {
	return &StorageFactory{}
}

func (f *StorageFactory) CreateClient(config *model.StorageConfig) (client.StorageClient, error) {
	switch config.Type {
	case "oss":
		return client.NewOssClient(config.OssConfig)
	default:
		return nil, fmt.Errorf("unknown storage type: %s", config.Type)
	}
}

func (f *StorageFactory) CreateClientByType(storageType string, configMap map[string]interface{}) (client.StorageClient, error) {
	switch storageType {
	case "oss":
		config := &model.OssConfig{
			Endpoint:        configMap["endpoint"].(string),
			AccessKeyId:     configMap["access_key_id"].(string),
			AccessKeySecret: configMap["access_key_secret"].(string),
			Bucket:          configMap["bucket"].(string),
		}
		return client.NewOssClient(config)
	default:
		return nil, fmt.Errorf("unknown storage type: %s", storageType)
	}
}
