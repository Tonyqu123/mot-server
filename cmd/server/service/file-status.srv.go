package service

import (
	// "encoding/json"
	"fmt"
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

func UpdateStatusByFileId(filestatus model.FileStatus) error {
	if err := model.UpdateStatusByFileId(filestatus); err != nil {
		return err
	}
	return nil
}

func DeleteStatusByFileId(fileId int) error {
	fmt.Println("DeleteStatusByFileId(fileId int)：", fileId)
	if err := model.DeleteStatusByFileId(fileId); err != nil {
		return err
	}
	return nil
}
