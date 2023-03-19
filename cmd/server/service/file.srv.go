package service

import (
	// "encoding/json"
	"github.com/tony/mot-server/cmd/server/model"
)

type FileSrv struct{}

func (a FileSrv) GetFiles() ([]model.File, int64, error) {
	files, err := model.GetFilesAndStatus()
	if err != nil {
		return nil, model.CountFiles(), err
	}

	return files, model.CountFiles(), nil
}

func AddFile(file model.File) (int, error) {
	_, err := model.AddFile(file)
	if err != nil {
		return -1, err
	}
	return 0, nil
}
