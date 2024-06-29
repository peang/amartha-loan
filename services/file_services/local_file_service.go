package file_services

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type localFileService struct{}

func NewLocalFileService() FileServiceInterface {
	return &localFileService{}
}

func (s *localFileService) Upload(file *multipart.FileHeader) (string, error) {
	uploadedFile, err := file.Open()
	if err != nil {
		return "", err
	}
	defer uploadedFile.Close()

	folderPath := "file_uploads"
	err = os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return "", err
	}

	filename := filepath.Join(folderPath, file.Filename)
	newFile, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer newFile.Close()

	// Copy the uploaded file to the new file
	_, err = io.Copy(newFile, uploadedFile)
	if err != nil {
		return "", err
	}

	return filename, nil
}
