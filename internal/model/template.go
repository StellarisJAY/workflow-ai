package model

import (
	"github.com/StellrisJAY/workflow-ai/internal/common"
	"time"
)

type Template struct {
	Id          int64     `json:"id,string" gorm:"primary_key;type:bigint"`
	Name        string    `json:"name" biding:"required" gorm:"column:name;type:varchar(50);not null"`
	Description string    `json:"description" gorm:"column:description;type:varchar(255);not null"`
	Data        string    `json:"data" biding:"required" gorm:"column:data;type:text;not null"`
	AddTime     time.Time `json:"addTime" gorm:"column:add_time;type:datetime;not null"`
	AddUser     int64     `json:"addUser" gorm:"column:add_user;type:bigint;not null"`
}

func (Template) TableName() string {
	return "wf_template"
}

type TemplateDetailDTO struct {
	Id          int64     `json:"id,string"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Data        string    `json:"data"`
	AddTime     time.Time `json:"addTime"`
	AddUser     int64     `json:"addUser"`
	AddUserName string    `json:"addUserName"`
	UsageCount  int64     `json:"usageCount"`
}

type TemplateListDTO struct {
	Id          int64     `json:"id,string"`
	Name        string    `json:"name"`
	AddTime     time.Time `json:"addTime"`
	AddUser     int64     `json:"addUser"`
	AddUserName string    `json:"addUserName"`
}

type TemplateQuery struct {
	Name string `json:"name"`
	common.PageQuery
}
