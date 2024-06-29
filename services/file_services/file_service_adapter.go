package file_services

import "mime/multipart"

type FileServiceInterface interface {
	Upload(file *multipart.FileHeader) (fileUrl string, err error)
}

func NewFileService() FileServiceInterface {
	// you can pass config here, and return instance depends on the config

	// for this purpose we're gonna local upload file service in file
	return NewLocalFileService()
}
