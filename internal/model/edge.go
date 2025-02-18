package model

type Edge struct {
	Id           string `json:"id"`
	Source       string `json:"source"`
	Target       string `json:"target"`
	Type         string `json:"type"`
	SourceHandle string `json:"source_handle"`
	TargetHandle string `json:"target_handle"`
}

type ConditionEdgeData struct {
}
