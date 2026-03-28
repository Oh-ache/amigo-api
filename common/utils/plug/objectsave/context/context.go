package context

import (
	"amigo-api/common/utils/plug/objectsave/client"
	"amigo-api/common/utils/plug/objectsave/model"
)

type StorageContext struct {
	client client.StorageClient
}

func NewStorageContext(client client.StorageClient) *StorageContext {
	return &StorageContext{
		client: client,
	}
}

func (s *StorageContext) SetClient(client client.StorageClient) {
	s.client = client
}

func (s *StorageContext) GetStorageType() string {
	return s.client.GetStorageType()
}

func (s *StorageContext) UploadFile(fileName string, data []byte) (*model.UploadResult, error) {
	return s.client.UploadFile(fileName, data)
}

func (s *StorageContext) UploadUrl(fileName string, url string) (*model.UploadResult, error) {
	return s.client.UploadUrl(fileName, url)
}

func (s *StorageContext) GetUploadToken(fileName string) (*model.UploadToken, error) {
	return s.client.GetUploadToken(fileName)
}
