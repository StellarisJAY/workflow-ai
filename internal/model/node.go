package model

type Node struct {
	Id       string `json:"id"`   // 节点ID
	Type     string `json:"type"` // 节点类型
	Position struct {
		X int `json:"x"`
		Y int `json:"y"`
	} `json:"position"` // 节点位置
	Data any `json:"data"` // 节点数据，根据节点类型不同，数据结构也不同
}

// LLMNodeData LLM节点数据
type LLMNodeData struct {
	Model           string   `json:"model"`           // 模型名称
	SystemPrompt    string   `json:"systemPrompt"`    // 系统提示词
	UserPrompt      string   `json:"userPrompt"`      // 用户提示词
	InputVariables  []string `json:"inputVariables"`  // 输入变量列表
	OutputFormat    string   `json:"outputFormat"`    // 输出格式 text,markdown,json
	OutputVariables []string `json:"outputVariables"` // 输出变量列表
}

// KnowledgeBaseWriteNodeData 写入知识库节点数据
type KnowledgeBaseWriteNodeData struct {
	KnowledgeBaseId int64  `json:"knowledgeBaseId"` // 知识库ID
	Content         []byte `json:"content"`         // 写入内容, {{nodeId.xxx}}表示写入某节点的变量列表的内容
}

// RetrieveKnowledgeBaseNodeData 检索知识库节点数据
type RetrieveKnowledgeBaseNodeData struct {
	KnowledgeBaseId int64  `json:"knowledgeBaseId"` // 知识库ID
	Query           string `json:"query"`           // 检索内容, {{nodeId.xxx}}表示从某节点的变量列表获取
}

type StartNodeData struct {
	InputVariables []string `json:"inputVariables"` // 输入变量列表
}

type EndNodeData struct {
	OutputVariables map[string]string `json:"outputVariables"` // 输出变量列表 key: 变量名，value: 变量来源 {{nodeId.xxx}}
}
