package model

import (
	"gorm.io/gorm"
)

type File struct {
	Filename string `json:"filename"`
	Fileid string `json:"fileid"`
	FileOrigin int `json:"file-origin"`
	FileTracked string `json:"file-tracked"`
	Userid string `json:"userid"`
	Uploadtime string `json:"upload-time"`
}

func (File) TableName() string {
	return "files"
}

func GetFiles() ([]File, error) {
	var files []File
	err := db.Find(&files).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return files, nil
}

func CountFiles() int64 {
	var total int64 = 0
	if err := db.Model(File{}).Count(&total).Error; err != nil {
		return -1
	}
	return total
}