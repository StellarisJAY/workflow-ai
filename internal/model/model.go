package model

//type ApiType string
//
//const (
//	ApiTypeOpenAI ApiType = "openai"
//	ApiTypeOllama ApiType = "ollama"
//)
//
//type Model struct {
//	Id        int64     `json:"id,string" gorm:"primary_key;column:id;type:bigint"`
//	Name      string    `json:"name" biding:"required" gorm:"column:name;type:varchar(50);not null"`
//	ApiType   string    `json:"apiType" biding:"required" gorm:"column:api_type;type:varchar(16);not null"`
//	BaseUrl   string    `json:"baseUrl" biding:"required" gorm:"column:base_url;type:varchar(255);not null"`
//	ApiKey    string    `json:"apiKey" gorm:"column:api_key;type:varchar(255)"`
//	Code      string    `json:"code" biding:"required" gorm:"column:code;type:varchar(32);not null"`
//	AddUser   int64     `json:"addUser" gorm:"column:add_user;type:bigint;not null"`
//	AddTime   time.Time `json:"addTime" gorm:"column:add_time;type:datetime;not null"`
//	ProviderModelType ProviderModelType `json:"modelType" biding:"required" gorm:"column:model_type;type:varchar(32);not null"`
//}
//
//func (Model) TableName() string {
//	return "wf_model"
//}

//type ModelDetailDTO struct {
//	Id          int64     `json:"id,string"`
//	Name        string    `json:"name"`
//	ApiType     string    `json:"apiType"`
//	BaseUrl     string    `json:"baseUrl"`
//	ApiKey      string    `json:"apiKey"`
//	Code        string    `json:"code"`
//	AddUser     int64     `json:"addUser"`
//	AddTime     time.Time `json:"addTime"`
//	AddUsername string    `json:"addUsername"`
//	ProviderModelType   ProviderModelType `json:"modelType"`
//}
//
//type ModelQuery struct {
//	Name      string    `json:"name" form:"name"`
//	ApiType   string    `json:"apiType" form:"apiType"`
//	Code      string    `json:"code" form:"code"`
//	ProviderModelType ProviderModelType `json:"modelType" form:"modelType"`
//	common.PageQuery
//}
//
//type ModelListDTO struct {
//	Id          int64     `json:"id,string"`
//	Name        string    `json:"name"`
//	ApiType     string    `json:"apiType"`
//	Code        string    `json:"code"`
//	AddUser     int64     `json:"addUser"`
//	AddTime     time.Time `json:"addTime"`
//	AddUsername string    `json:"addUsername"`
//	ProviderModelType   ProviderModelType `json:"modelType"`
//}
