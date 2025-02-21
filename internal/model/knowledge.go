package model

import (
	"encoding/json"
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
	KbFileProcessing
	KbFileProcessed
	KbFileFailed
)

func (s KbFileStatus) String() string {
	switch s {
	case KbFileStatusUnavailable:
		return "无效"
	case KbFileUploaded:
		return "未解析"
	case KbFileProcessing:
		return "解析中"
	case KbFileProcessed:
		return "解析完成"
	case KbFileFailed:
		return "解析失败"
	}
	return "无效"
}

type KnowledgeBaseFile struct {
	Id       int64        `json:"id,string" gorm:"primary_key;column:id;type:bigint;"`
	Name     string       `json:"name" binding:"required" gorm:"column:name;type:varchar(50);not null"`
	KbId     int64        `json:"kbId,string" binding:"required" gorm:"column:kb_id;type:bigint;not null"`
	AddTime  time.Time    `json:"addTime" gorm:"column:add_time;type:datetime;not null"`
	AddUser  int64        `json:"addUser,string" gorm:"column:add_user;type:bigint;not null"`
	Length   int64        `json:"length,string" gorm:"column:length;type:bigint;not null"` // 文件大小
	Chunks   int64        `json:"chunks,string" gorm:"column:chunks;type:bigint;not null"` // 分片数量
	Url      string       `json:"url" gorm:"column:url;type:varchar(255);not null"`        // 文件地址
	Status   KbFileStatus `json:"status" gorm:"column:status;type:tinyint;not null"`
	Metadata string       `json:"metadata" gorm:"column:metadata;type:text;not null"` // 文件元数据 kv
}

func (KnowledgeBaseFile) TableName() string {
	return "wf_knowledge_base_file"
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

type KbFileProcessStatus int

const (
	KbFileProcessStatusQueued KbFileProcessStatus = iota
	KbFileProcessStatusSplitting
	KbFileProcessStatusEmbedding
	KbFileProcessStatusCompleted
	KbFileProcessStatusFailed
)

type KbFileProcessRequest struct {
}

type KbFileProcessTask struct {
	Id           int64               `json:"id,string" gorm:"primary_key;column:id;type:bigint;"`
	KbId         int64               `json:"kbId,string" gorm:"column:kb_id;type:bigint;not null;uniqueIndex:kb_file_idx;"`
	FileId       int64               `json:"fileId,string" gorm:"column:file_id;type:bigint;not null;uniqueIndex:kb_file_idx;"`
	Status       KbFileProcessStatus `json:"status" gorm:"column:status;type:tinyint;not null;"`
	AddTime      time.Time           `json:"addTime" gorm:"column:add_time;type:datetime;not null"`
	CompleteTime time.Time           `json:"completeTime" gorm:"column:complete_time;type:datetime;not null"`
	Error        string              `json:"error" gorm:"column:error;type:text;not null"`
}

func (KbFileProcessTask) TableName() string {
	return "wf_kb_file_process_task"
}

// KbFileProcessOptions 知识库文件解析选项
type KbFileProcessOptions struct {
	FileId     int64  `json:"fileId,string" gorm:"primary key;column:file_id;type:bigint;not null"` // 文件id
	Separators string `json:"separators" gorm:"column:separators;type:varchar(255);not null"`       // 分隔符
	ChunkSize  int    `json:"chunkSize" gorm:"column:chunk_size;type:int;not null"`                 // 分片大小
}

type KbFileProcessOptionsUpdateDTO struct {
	FileId     int64    `json:"fileId,string"` // 文件id
	Separators []string `json:"separators"`    // 分隔符
	ChunkSize  int      `json:"chunkSize"`     // 分片大小
}

func (KbFileProcessOptions) TableName() string {
	return "wf_kb_file_process_options"
}

var defaultSeparators = []string{
	" ",
	".",
	",",
	"\u200b", // Zero-width space
	"\uff0c", // Fullwidth comma
	"\u3001", // Ideographic comma
	"\uff0e", // Fullwidth full stop
	"\u3002", // Ideographic full stop
	"",
	"\n",
	"！",
	"？",
	"!",
	"?",
}

func DefaultKbFileProcessOptions() KbFileProcessOptions {
	data, _ := json.Marshal(defaultSeparators)
	return KbFileProcessOptions{
		Separators: string(data),
		ChunkSize:  512,
	}
}

type KbFileChunk struct {
	Id      int64 `json:"id,string" gorm:"primary_key;column:id;type:bigint;"`
	FileId  int64 `json:"fileId,string" gorm:"column:file_id;type:bigint;not null"`
	ChunkId int64 `json:"chunkId,string" gorm:"column:chunk_id;type:bigint;not null"`
}

func (KbFileChunk) TableName() string {
	return "wf_kb_file_chunk"
}

type KbSearchRequest struct {
	KbId           int64   `json:"kbId,string"`
	Input          string  `json:"input"`
	ScoreThreshold float32 `json:"scoreThreshold"`
	Count          int     `json:"count"`
}

type KbSearchReturnDocument struct {
	Content string  `json:"content"`
	ChunkId string  `json:"chunkId"`
	Score   float32 `json:"score,string"`
	FileId  string  `json:"fileId"`
}

type KbSearchResult struct {
	Documents []*KbSearchReturnDocument `json:"documents"`
	Files     []*KbFileListDTO          `json:"files"`
}
