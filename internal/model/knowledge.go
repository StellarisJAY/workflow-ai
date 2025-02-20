package model

import (
	"github.com/StellrisJAY/workflow-ai/internal/common"
	"time"
)

type KnowledgeBase struct {
	Id             int64     `json:"id,string" gorm:"primary_key;column:id;type:bigint"`
	Name           string    `json:"name" binding:"required" gorm:"column:name;type:varchar(50);not null"`
	AddTime        time.Time `json:"addTime" gorm:"column:add_time;type:datetime;not null"`
	AddUser        int64     `json:"addUser,string" gorm:"column:add_user;type:bigint;not null"`
	Enable         bool      `json:"enable" gorm:"column:enable;type:boolean;not null"`
	EmbeddingModel int64     `json:"embeddingModel,string" binding:"required" gorm:"column:embedding_model;type:bigint;not null"` // 知识库使用的嵌入模型
	Description    string    `json:"description" binding:"required" gorm:"column:description;not null;type:varchar(255)"`
}

func (KnowledgeBase) TableName() string {
	return "wf_knowledge_base"
}

type KbFileStatus int

const (
	KbFileStatusUnavailable KbFileStatus = iota
	KbFileUploaded
	KbFileProcessed
	KbFileFailed
)

func (s KbFileStatus) String() string {
	switch s {
	case KbFileStatusUnavailable:
		return "无效"
	case KbFileUploaded:
		return "未解析"
	case KbFileProcessed:
		return "解析完成"
	case KbFileFailed:
		return "解析失败"
	}
	return "无效"
}

type KnowledgeBaseFile struct {
	Id      int64        `json:"id,string" gorm:"primary_key;column:id;type:bigint;"`
	Name    string       `json:"name" binding:"required" gorm:"column:name;type:varchar(50);not null"`
	KbId    int64        `json:"kbId,string" binding:"required" gorm:"column:kb_id;type:bigint;not null"`
	AddTime time.Time    `json:"addTime" gorm:"column:add_time;type:datetime;not null"`
	AddUser int64        `json:"addUser,string" gorm:"column:add_user;type:bigint;not null"`
	Length  int64        `json:"length,string" gorm:"column:length;type:bigint;not null"` // 文件大小
	Chunks  int64        `json:"chunks,string" gorm:"column:chunks;type:bigint;not null"` // 分片数量
	Url     string       `json:"url" gorm:"column:url;type:varchar(255);not null"`        // 文件地址
	Status  KbFileStatus `json:"status" gorm:"column:status;type:tinyint;not null"`
}

func (KnowledgeBaseFile) TableName() string {
	return "wf_knowledge_base_file"
}

type KnowledgeBaseImportTask struct {
	Id       int64  `json:"id,string"`
	KbId     int64  `json:"kbId,string"`
	KbFileId int64  `json:"kbFileId,string"`
	Status   string `json:"status"`
	AddTime  string `json:"addTime"`
}

type KnowledgeBaseQuery struct {
	common.PageQuery
}

type KnowledgeBaseListDTO struct {
	Id            int64     `json:"id,string"`
	Name          string    `json:"name"`
	AddTime       time.Time `json:"addTime"`
	AddUser       int64     `json:"addUser,string"`
	Enable        bool      `json:"Enable"`
	Description   string    `json:"description"`
	DocumentCount int       `json:"documentCount"`
	Size          int64     `json:"size"`
}

type KnowledgeBaseDetailDTO struct {
	Id             int64     `json:"id,string"`
	Name           string    `json:"name"`
	AddTime        time.Time `json:"addTime"`
	AddUser        int64     `json:"addUser,string"`
	Enable         bool      `json:"enable"`
	EmbeddingModel int64     `json:"embeddingModel,string"` // 知识库使用的嵌入模型
	Description    string    `json:"description"`
	DocumentCount  int       `json:"documentCount"`
	Size           int64     `json:"size"`
}

type KbFileCountSize struct {
	Count int64 `gorm:"column:count" json:"count"`
	Size  int64 `gorm:"column:size" json:"size"`
	KbId  int64 `gorm:"column:kb_id" json:"kb_id"`
}

type KbFileQuery struct {
	common.PageQuery
}

type KbFileListDTO struct {
	Id         int64        `json:"id,string"`
	Name       string       `json:"name"`
	AddTime    time.Time    `json:"addTime"`
	AddUser    int64        `json:"addUser,string"`
	Length     int64        `json:"length"`
	Status     KbFileStatus `json:"status"`
	StatusName string       `json:"statusName"`
}
