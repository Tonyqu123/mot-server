package model

import (
	// "fmt"
	"gorm.io/gorm"
)

type File struct {
	gorm.Model

	// foreign key, ref: https://gorm.io/docs/has_many.html
	UserID      uint
	Filename    string
	FileOrigin  string
	FileTracked string
	UploadTime  string

	FileStatus FileStatus
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

func AddFile(file File) (int, error) {
	var err error

	sql := db.Create(&file)
	if err = sql.Error; err != nil {
		return -1, err
	}
	return 0, nil
}
