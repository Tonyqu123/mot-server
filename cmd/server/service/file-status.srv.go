package service

import (
	// "encoding/json"
	"github.com/tony/mot-server/cmd/server/model"
)

type FileStatusSrv struct{}

func (a FileStatusSrv) GetFileStatus() ([]model.File, int64, error) {
	files, err := model.GetFilesAndStatus()
	if err != nil {
		return nil, model.CountFiles(), err
	}

	return files, model.CountFiles(), nil
}

func AddFileStatus(filestatus model.FileStatus) error {
	if err := model.AddFileStatus(filestatus); err != nil {
		return err
	}
	return nil
}
