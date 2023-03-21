package service

import (
	// "encoding/json"
	"fmt"
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

func AddFile(file model.File) error {
	fileId, err := model.AddFile(file)
	var fileStatus model.FileStatus
	fileStatus.FileID = uint(fileId)
	fileStatus.Status = 0
	fmt.Println("model.AddFileStatus：", fileStatus)
	err = model.AddFileStatus(fileStatus)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	return nil
}

func (a FileSrv) DeleteFileById(fileId int) error {
	fmt.Println("(a FileSrv) DeleteFileById：", fileId)
	// 同时删除 file 和 filestatus 数据库的内容
	err := model.DeleteFileById(fileId)
	if err != nil {
		return err
	}
	err = DeleteStatusByFileId(fileId)
	if err != nil {
		return err
	}
	return nil
}
