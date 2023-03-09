package model

import (
	"gorm.io/gorm"
)

type Filestatus struct {
	Fileid  int `json:"fileid"`
	Status  int `json:"status"`
}

type FileAndStatus struct {
	Filename string `json:"filename"`
	Fileid int `json:"fileid" gorm:"autoIncrement"`
	FileOrigin string `json:"file_origin" gorm:"column:file_origin"`
	FileTracked string `json:"file_tracked" gorm:"column:file_tracked"`
	Userid int `json:"userid"`
	Uploadtime string `json:"upload_time" gorm:"column:upload_time"`
	Status int `json:"status" gorm:"status"`
}

func (Filestatus) TableName() string {
	return "filestatus"
}

func GetFilesAndStatus() ([]FileAndStatus, error) {
	var files []FileAndStatus
	// err := db.Find(&files).Error
	err := db.Model(&File{}).Select("files.filename, files.fileid, files.file_origin, files.file_tracked, files.userid, files.upload_time, filestatus.status").Joins("inner join filestatus on filestatus.fileid = files.fileid").Scan(&files).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return files, nil
}