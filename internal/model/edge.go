package model

type Edge struct {
	Id     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
	Type   string `json:"type"`
	Data   any    `json:"data"`
}

type ConditionEdgeData struct {
}
