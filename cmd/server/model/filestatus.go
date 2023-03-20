package model

import (
	"gorm.io/gorm"
)

type FileStatus struct {
	// foreign key, ref:https://gorm.io/docs/has_one.html
	gorm.Model
	FileID uint `json:"file_id"`
	Status int  `json:"status"`
}

//type FileAndStatus struct {
//	Filename    string `json:"filename"`
//	Fileid      int    `json:"file_id"`
//	FileOrigin  string `json:"file_origin" gorm:"column:file_origin"`
//	FileTracked string `json:"file_tracked" gorm:"column:file_tracked"`
//	UserID      int    `json:"user_id"`
//	Uploadtime  string `json:"upload_time" gorm:"column:upload_time"`
//	Status      int    `json:"status" gorm:"status"`
//}

func (FileStatus) TableName() string {
	return "filestatus"
}

func GetFilesAndStatus() ([]File, error) {
	var files []File
	var err error

	// ref:
	err = db.Model(&File{}).Preload("FileStatus").Find(&files).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return files, nil
}

func AddFileStatus(fileStatus FileStatus) error {
	if err := db.Create(&fileStatus).Error; err != nil {
		return err
	}
	return nil
}

func (FileStatus) UpdateStatusByFileId(fileId int, statusId int) {
	db.Model(&FileStatus{}).Where("file_id = ?", fileId).Update("status", statusId)
}
