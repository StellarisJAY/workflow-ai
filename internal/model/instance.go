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
	Id           int64                  `json:"id,string"`
	TemplateId   int64                  `json:"templateId,string"`
	Data         string                 `json:"data"`
	Status       WorkflowInstanceStatus `json:"status"`
	AddTime      time.Time              `json:"addTime"`
	AddUser      int64                  `json:"addUser"`
	CompleteTime time.Time              `json:"completeTime"`
}

// NodeInstance 节点实例表
type NodeInstance struct {
	Id           int64              `json:"id"`
	WorkflowId   int64              `json:"workflowId"`
	Type         string             `json:"type"`
	NodeId       string             `json:"nodeId"`
	AddTime      time.Time          `json:"addTime"`
	CompleteTime time.Time          `json:"completeTime"`
	Status       NodeInstanceStatus `json:"status"`
	Output       string             `json:"output"` // 节点输出变量json
	Error        string             `json:"error"`  // 节点执行错误信息
}

type NodeInstanceOutputDTO struct {
	Id           int64              `json:"id"`
	NodeId       string             `json:"nodeId"`
	NodeName     string             `json:"nodeName"`
	Type         string             `json:"type"`
	AddTime      time.Time          `json:"addTime"`
	CompleteTime time.Time          `json:"completeTime"`
	Status       NodeInstanceStatus `json:"status"`
	Output       string             `json:"output"`
	Error        string             `json:"error"`
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
	Id             int64                  `json:"id,string"`
	TemplateId     int64                  `json:"templateId,string"`
	TemplateName   string                 `json:"templateName"`
	Status         WorkflowInstanceStatus `json:"status"`
	StatusName     string                 `json:"statusName"`
	AddTime        time.Time              `json:"addTime"`
	AddUser        int64                  `json:"addUser"`
	CompleteTime   time.Time              `json:"completeTime"`
	Duration       string                 `json:"duration"`
	Data           string                 `json:"data"`
	NodeStatusList []*NodeStatusDTO       `json:"nodeStatusList" gorm:"-"`
}

type NodeStatusDTO struct {
	Id         int64              `json:"id"`
	NodeId     string             `json:"nodeId"`
	Status     NodeInstanceStatus `json:"status"`
	StatusName string             `json:"statusName"`
}

type NodeInstanceDetailDTO struct {
	Id           int64              `json:"id"`
	WorkflowId   int64              `json:"workflowId"`
	Type         string             `json:"type"`
	NodeId       string             `json:"nodeId"`
	AddTime      time.Time          `json:"addTime"`
	CompleteTime time.Time          `json:"completeTime"`
	Status       NodeInstanceStatus `json:"status"`
	Output       string             `json:"output"` // 节点输出变量json
	Error        string             `json:"error"`  // 节点执行错误信息
	StatusName   string             `json:"statusName"`
}
