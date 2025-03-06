package model

import "time"

type WorkflowInstanceStatus int

const (
	WorkflowInstanceStatusRunning WorkflowInstanceStatus = iota
	WorkflowInstanceStatusCompleted
	WorkflowInstanceStatusFailed
)

func (wi WorkflowInstanceStatus) String() string {
	switch wi {
	case WorkflowInstanceStatusRunning:
		return "运行中"
	case WorkflowInstanceStatusCompleted:
		return "完成"
	case WorkflowInstanceStatusFailed:
		return "失败"
	default:
		return "Unknown"
	}
}

type NodeInstanceStatus int

const (
	NodeInstanceStatusUnreached NodeInstanceStatus = iota
	NodeInstanceStatusRunning
	NodeInstanceStatusCompleted
	NodeInstanceStatusFailed
)

func (ni NodeInstanceStatus) String() string {
	switch ni {
	case NodeInstanceStatusUnreached:
		return "未到达"
	case NodeInstanceStatusRunning:
		return "运行中"
	case NodeInstanceStatusCompleted:
		return "完成"
	case NodeInstanceStatusFailed:
		return "失败"
	default:
		return "Unknown"
	}
}

type WorkflowInstance struct {
	Id           int64                  `json:"id,string" gorm:"primary_key;column:id;type:bigint"`
	TemplateId   int64                  `json:"templateId,string" gorm:"column:template_id;type:bigint;not null"`
	Data         string                 `json:"data" gorm:"column:data;type:longtext;not null"`
	Status       WorkflowInstanceStatus `json:"status" gorm:"column:status;type:int;not null"`
	AddTime      time.Time              `json:"addTime" gorm:"column:add_time;type:datetime;not null"`
	AddUser      int64                  `json:"addUser" gorm:"column:add_user;type:bigint;not null"`
	CompleteTime time.Time              `json:"completeTime" gorm:"column:complete_time;type:datetime;not null"`
}

func (WorkflowInstance) TableName() string {
	return "wf_workflow_instance"
}

// NodeInstance 节点实例表
type NodeInstance struct {
	Id           int64              `json:"id" gorm:"primary_key;column:id;type:bigint"`
	WorkflowId   int64              `json:"workflowId" gorm:"column:workflow_id;type:bigint;not null"`
	Type         string             `json:"type" gorm:"column:type;type:varchar(32);not null"`
	NodeId       string             `json:"nodeId" gorm:"column:node_id;type:varchar(64);not null"`
	AddTime      time.Time          `json:"addTime" gorm:"column:add_time;type:datetime;not null"`
	CompleteTime time.Time          `json:"completeTime" gorm:"column:complete_time;type:datetime;not null"`
	Status       NodeInstanceStatus `json:"status" gorm:"column:status;type:int;not null"`
	Output       string             `json:"output" gorm:"column:output;type:json"` // 节点输出变量json
	Error        string             `json:"error" gorm:"column:error;type:text"`   // 节点执行错误信息
}

func (NodeInstance) TableName() string {
	return "wf_node_instance"
}

type WorkflowInstanceListDTO struct {
	Id           int64                  `json:"id,string"`
	TemplateId   int64                  `json:"templateId,string"`
	TemplateName string                 `json:"templateName"`
	Status       WorkflowInstanceStatus `json:"status"`
	StatusName   string                 `json:"statusName"`
	AddTime      time.Time              `json:"addTime"`
	AddUser      int64                  `json:"addUser"`
	CompleteTime time.Time              `json:"completeTime"`
	Duration     string                 `json:"duration"`
}

type WorkflowInstanceDetailDTO struct {
	Id                int64                               `json:"id,string"`
	TemplateId        int64                               `json:"templateId,string"`
	TemplateName      string                              `json:"templateName"`
	Status            WorkflowInstanceStatus              `json:"status"`
	StatusName        string                              `json:"statusName"`
	AddTime           time.Time                           `json:"addTime"`
	AddUser           int64                               `json:"addUser"`
	CompleteTime      time.Time                           `json:"completeTime"`
	Duration          string                              `json:"duration"`
	Data              string                              `json:"data"`
	NodeStatusList    []*NodeStatusDTO                    `json:"nodeStatusList" gorm:"-"`
	PassedEdgesList   []string                            `json:"passedEdgesList" gorm:"-"`
	SuccessBranchList []*WorkflowInstanceSuccessBranchDTO `json:"successBranchList" gorm:"-"`
}

type WorkflowInstanceSuccessBranchDTO struct {
	NodeId string `json:"nodeId"`
	Branch string `json:"branch"`
}

type NodeStatusDTO struct {
	Id         int64              `json:"id"`
	NodeId     string             `json:"nodeId"`
	Status     NodeInstanceStatus `json:"status"`
	StatusName string             `json:"statusName"`
}

type NodeInstanceDetailDTO struct {
	Id                  int64                   `json:"id"`
	WorkflowId          int64                   `json:"workflowId"`
	Type                string                  `json:"type"`
	NodeId              string                  `json:"nodeId"`
	AddTime             time.Time               `json:"addTime"`
	CompleteTime        time.Time               `json:"completeTime"`
	Status              NodeInstanceStatus      `json:"status"`
	Output              string                  `json:"output"` // 节点输出变量json
	Error               string                  `json:"error"`  // 节点执行错误信息
	StatusName          string                  `json:"statusName"`
	OutputVariableTypes map[string]VariableType `json:"outputVariableTypes" gorm:"-"`
}

type WorkflowInstanceTimelineDTO struct {
	Id           int64              `json:"id,string"`
	NodeId       string             `json:"nodeId"`
	NodeName     string             `json:"nodeName"`
	Status       NodeInstanceStatus `json:"status"`
	StatusName   string             `json:"statusName"`
	AddTime      time.Time          `json:"addTime"`
	CompleteTime time.Time          `json:"completeTime"`
	Duration     string             `json:"duration"`
}

type NodeInstanceOutputVariable struct {
	Type  VariableType `json:"type"`
	Value any          `json:"value"`
}
