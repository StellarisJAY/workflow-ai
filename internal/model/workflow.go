package model

import "time"

// Workflow 工作流表
type Workflow struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	AddTime     time.Time `json:"addTime"`
	AddUser     int64     `json:"addUser"`
	Nodes       string    `json:"nodes"`
	Edges       string    `json:"edges"`
}

// WorkflowDefinition 工作流定义json
type WorkflowDefinition struct {
	Nodes []*Node `json:"nodes"`
	Edges []*Edge `json:"edges"`
}
