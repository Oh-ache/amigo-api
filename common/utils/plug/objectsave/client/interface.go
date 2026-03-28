package client

import (
	"amigo-api/common/utils/plug/objectsave/model"
)

type StorageClient interface {
	UploadFile(fileName string, data []byte) (*model.UploadResult, error)
	UploadUrl(fileName string, url string) (*model.UploadResult, error)
	GetUploadToken(fileName string) (*model.UploadToken, error)
	GetStorageType() string
}
