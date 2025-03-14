package model

// WorkflowDefinition 工作流定义json
type WorkflowDefinition struct {
	Nodes []*Node `json:"nodes"`
	Edges []*Edge `json:"edges"`
}

type StartWorkflowRequest struct {
	TemplateId int64          `json:"templateId,string" biding:"required"`
	Inputs     map[string]any `json:"inputs" biding:"required"`
	Definition string         `json:"definition"`
}
