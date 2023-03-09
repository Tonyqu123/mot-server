package service

import (
	// "encoding/json"
	"github.com/tony/mot-server/cmd/server/model"
)

type FileSrv struct{}

func (a FileSrv) GetFiles() ([]model.FileAndStatus, int64, error) {
	files, err := model.GetFilesAndStatus()
	total := model.CountFiles()
	if err != nil {
		return nil, total, err
	}
	return files, total, nil
}