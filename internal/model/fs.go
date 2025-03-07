package model

import "time"

type File struct {
	Id      int64     `json:"id" gorm:"primary_key;type:bigint;"`
	Name    string    `json:"name" gorm:"type:varchar(255);not null"`
	Type    string    `json:"type" gorm:"type:varchar(32);not null"`
	Size    int64     `json:"size" gorm:"type:bigint;not null"`
	Url     string    `json:"url" gorm:"type:varchar(255);not null"`
	AddTime time.Time `json:"add_time" gorm:"type:datetime;not null"`
	AddUser int64     `json:"add_user" gorm:"type:bigint;not null"`
	Md5     string    `json:"md5" gorm:"type:varchar(32);not null"`
}

func (File) TableName() string {
	return "wf_file"
}
