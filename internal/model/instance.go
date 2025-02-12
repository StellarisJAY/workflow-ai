package model

import "time"

type WorkflowInstanceStatus int

const (
	WorkflowInstanceStatusRunning WorkflowInstanceStatus = iota
	WorkflowInstanceStatusCompleted
	WorkflowInstanceStatusFailed
)

type NodeInstanceStatus int

const (
	NodeInstanceStatusUnreached NodeInstanceStatus = iota
	NodeInstanceStatusRunning
	NodeInstanceStatusCompleted
	NodeInstanceStatusFailed
)

type WorkflowInstance struct {
	Id           int64                  `json:"id"`
	TemplateId   int64                  `json:"templateId"`
	Data         string                 `json:"data"`
	Input        string                 `json:"input"`
	Output       string                 `json:"output"`
	Status       WorkflowInstanceStatus `json:"status"`
	AddTime      time.Time              `json:"addTime"`
	AddUser      int64                  `json:"addUser"`
	CompleteTime time.Time              `json:"completeTime"`
}

// NodeInstance 节点实例表
type NodeInstance struct {
	Id           int64              `json:"id"`
	WorkflowId   int64              `json:"workflowId"`
	NodeId       string             `json:"nodeId"`
	AddTime      time.Time          `json:"addTime"`
	CompleteTime time.Time          `json:"completeTime"`
	Status       NodeInstanceStatus `json:"status"`
	Input        string             `json:"input"`  // 节点输入变量json
	Output       string             `json:"output"` // 节点输出变量json
	Error        string             `json:"error"`  // 节点执行错误信息
}
