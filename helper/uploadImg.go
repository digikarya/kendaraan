package helper

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func ImgProcess(r *http.Request,path,filename string) (string,error) {
	if err := r.ParseMultipartForm(2048); err != nil {
		return "",err
	}
	uploadedFile, handler, err := r.FormFile("image")
	if err != nil {
		return "",err
	}
	defer uploadedFile.Close()
	dir, err := os.Getwd()
	if err != nil {
		return "",err
	}
	isContinue := true
	switch filepath.Ext(handler.Filename) {
	case ".png":
		isContinue = true
	case ".jpg":
		isContinue = true
	case ".jpeg":
		isContinue = true

	default:
		isContinue = false

	}
	if isContinue != true{
		return "", errors.New("Invalid file type")
	}
	filename =filename+filepath.Ext(handler.Filename)
	fileLocation := filepath.Join(dir, "../files/"+path, filename)
	targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "",err
	}
	defer targetFile.Close()
	if _, err := io.Copy(targetFile, uploadedFile); err != nil {
		return "",err
	}
	return filename,nil
}

