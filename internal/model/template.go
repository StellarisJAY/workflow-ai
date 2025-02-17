package model

import "time"

type Template struct {
	Id          int64     `json:"id,string"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Data        string    `json:"data"`
	AddTime     time.Time `json:"addTime"`
	AddUser     int64     `json:"addUser"`
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
}
