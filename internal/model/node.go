package model

type NodeType string

const (
	NodeTypeStart              NodeType = "start"              // 开始节点
	NodeTypeLLM                NodeType = "llm"                // 大模型节点
	NodeTypeKnowledgeRetrieval NodeType = "knowledgeRetrieval" // 知识库检索节点
	NodeTypeKnowledgeWrite     NodeType = "knowledgeWrite"     // 写入知识库节点
	NodeTypeEnd                NodeType = "end"                // 结束节点
	NodeTypeCrawler            NodeType = "crawler"            // 爬虫节点
)

type Node struct {
	Id       string `json:"id"`   // 节点ID
	Type     string `json:"type"` // 节点类型
	Position struct {
		X int `json:"x"`
		Y int `json:"y"`
	} `json:"position"` // 节点位置
	Data struct {
		LLMNodeData                   *LLMNodeData                   `json:"llmNodeData"`
		KnowledgeBaseWriteNodeData    *KnowledgeBaseWriteNodeData    `json:"knowledgeBaseWriteNodeData"`
		RetrieveKnowledgeBaseNodeData *RetrieveKnowledgeBaseNodeData `json:"retrieveKnowledgeBaseNodeData"`
		StartNodeData                 *StartNodeData                 `json:"startNodeData"`
		EndNodeData                   *EndNodeData                   `json:"endNodeData"`
	} `json:"data"`
}

// LLMNodeData LLM节点数据
type LLMNodeData struct {
	ModelName       string            `json:"model"`           // 模型名称
	ModelId         int64             `json:"modelId"`         // 模型ID
	SystemPrompt    string            `json:"systemPrompt"`    // 系统提示词
	UserPrompt      string            `json:"userPrompt"`      // 用户提示词
	InputVariables  map[string]string `json:"inputVariables"`  // 输入变量列表, key:变量名，value：变量来源{{nodeId.xxx}}或空(运行时输入)
	OutputFormat    string            `json:"outputFormat"`    // 输出格式 text,markdown,json
	OutputVariables []string          `json:"outputVariables"` // 输出变量名列表
}

// KnowledgeBaseWriteNodeData 写入知识库节点数据
type KnowledgeBaseWriteNodeData struct {
	KnowledgeBaseId int64  `json:"knowledgeBaseId"` // 知识库ID
	Content         string `json:"content"`         // 写入内容, {{nodeId.xxx}}表示写入某节点的变量列表的内容, 空表示运行时输入
}

// RetrieveKnowledgeBaseNodeData 检索知识库节点数据
type RetrieveKnowledgeBaseNodeData struct {
	KnowledgeBaseId int64  `json:"knowledgeBaseId"` // 知识库ID
	Query           string `json:"query"`           // 检索内容, {{nodeId.xxx}}表示从某节点的变量列表获取，空表示运行时输入
}

type StartNodeData struct {
	InputVariables []string `json:"inputVariables"` // 输入变量列表
}

type EndNodeData struct {
	OutputVariables map[string]string `json:"outputVariables"` // 输出变量列表 key: 变量名，value: 变量来源 {{nodeId.xxx}}
}
