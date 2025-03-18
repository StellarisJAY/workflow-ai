package model

import (
	"github.com/StellrisJAY/workflow-ai/internal/common"
	"time"
)

type ProviderModelType string

const (
	ProviderModelTypeLargeLanguageModel ProviderModelType = "llm"                 // llm
	ProviderModelTypeEmbedding          ProviderModelType = "embedding"           // 嵌入
	ProviderModelTypeImageUnderstanding ProviderModelType = "image_understanding" // 视觉语言
)

type ProviderCode string

const (
	ProviderCodeOpenAI ProviderCode = "openai"
	ProviderCodeOllama ProviderCode = "ollama"
	ProviderCodeTongyi ProviderCode = "tongyi"
)

// Provider 供应商
type Provider struct {
	Id          int64        `json:"id,string" gorm:"primary_key;column:id;type:bigint"`
	Name        string       `json:"name" biding:"required" gorm:"column:name;type:varchar(50);not null"`
	Code        ProviderCode `json:"code" biding:"required" gorm:"column:code;type:varchar(50);not null"`
	Credentials string       `json:"credentials" gorm:"column:credentials;type:text"`
	AddTime     time.Time    `json:"addTime" gorm:"column:add_time;type:datetime;not null"`
}

func (Provider) TableName() string {
	return "wf_provider"
}

type NewProviderSchema struct {
	Code             ProviderCode `json:"code"`
	Name             string       `json:"name"`
	CredentialSchema any          `json:"credentialSchema"`
}

// ProviderModel 供应商模型
type ProviderModel struct {
	Id          int64             `json:"id,string" gorm:"primary_key;column:id;type:bigint"`
	ProviderId  int64             `json:"providerId,string" gorm:"column:provider_id;type:bigint"`
	ModelName   string            `json:"modelName" gorm:"column:model_name;type:varchar(64);not null"`
	ModelType   ProviderModelType `json:"modelType" gorm:"column:model_type;type:varchar(32);not null"`
	Credentials string            `json:"credentials" gorm:"column:credentials;type:text"`
	AddTime     time.Time         `json:"addTime" gorm:"column:add_time;type:datetime;not null"`
}

func (ProviderModel) TableName() string {
	return "wf_provider_model"
}

type ProviderListDTO struct {
	Id      int64                   `json:"id,string"`
	Name    string                  `json:"name"`
	AddTime time.Time               `json:"addTime"`
	Code    ProviderCode            `json:"code"`
	Models  []*ProviderModelListDTO `json:"models" gorm:"-"`
}

type ProviderModelListDTO struct {
	Id           int64             `json:"id,string"`
	ProviderId   int64             `json:"providerId,string"`
	ProviderName string            `json:"providerName"`
	ModelName    string            `json:"modelName"`
	ModelType    ProviderModelType `json:"modelType"`
	AddTime      time.Time         `json:"addTime"`
	ProviderCode ProviderCode      `json:"providerCode"`
}

type ProviderModelDetail struct {
	Id                  int64             `json:"id,string"`
	ProviderId          int64             `json:"providerId,string"`
	ProviderName        string            `json:"providerName"`
	ModelName           string            `json:"modelName"`
	ModelType           ProviderModelType `json:"modelType"`
	AddTime             time.Time         `json:"addTime"`
	Credentials         string            `json:"credentials"`
	ProviderCredentials string            `json:"providerCredentials"`
	ProviderCode        ProviderCode      `json:"providerCode"`
}

type ProviderModelQuery struct {
	ProviderId int64             `json:"providerId,string" form:"providerId"`
	ModelType  ProviderModelType `json:"modelType" form:"modelType"`
	common.PageQuery
}

type OpenAICredentials struct {
	ApiKey  string `json:"apiKey"`
	BaseUrl string `json:"baseUrl"`
}

type TongyiCredentials struct {
	ApiKey string `json:"apiKey"`
}

var ProviderSchemas = []NewProviderSchema{
	{
		Code:             ProviderCodeOpenAI,
		Name:             "OpenAI",
		CredentialSchema: OpenAICredentials{},
	},
	{
		Code:             ProviderCodeOllama,
		Name:             "Ollama",
		CredentialSchema: struct{}{},
	},
	{
		Code:             ProviderCodeTongyi,
		Name:             "通义千问",
		CredentialSchema: TongyiCredentials{},
	},
}
