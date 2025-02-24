package model

type NodeType string

const (
	NodeTypeStart              NodeType = "start"              // 开始节点
	NodeTypeLLM                NodeType = "llm"                // 大模型节点
	NodeTypeKnowledgeRetrieval NodeType = "knowledgeRetrieval" // 知识库检索节点
	NodeTypeKnowledgeWrite     NodeType = "knowledgeWrite"     // 写入知识库节点
	NodeTypeEnd                NodeType = "end"                // 结束节点
	NodeTypeCrawler            NodeType = "crawler"            // 爬虫节点
	NodeTypeCondition          NodeType = "condition"
)

type VariableType string

const (
	VariableTypeNumber VariableType = "number"
	VariableTypeString VariableType = "string"
	VariableTypeFile   VariableType = "file"
	VariableTypeRef    VariableType = "ref" // 引用其他节点的变量, 值为 节点ID.变量名
	VariableTypeArray  VariableType = "array"
)

type KbSearchType string

const (
	KbSearchTypeSimilarity KbSearchType = "similarity" // 语义搜索
	KbSearchTypeFulltext   KbSearchType = "fulltext"   // 全文搜索
	KbSearchTypeMixed      KbSearchType = "mixed"      // 混合搜索
)

type Node struct {
	Id       string `json:"id"`   // 节点ID
	Type     string `json:"type"` // 节点类型
	Position struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	} `json:"position"` // 节点位置
	Data struct {
		Name                          string                         `json:"name"`
		LLMNodeData                   *LLMNodeData                   `json:"llmNodeData"`
		KnowledgeBaseWriteNodeData    *KnowledgeBaseWriteNodeData    `json:"knowledgeBaseWriteNodeData"`
		RetrieveKnowledgeBaseNodeData *RetrieveKnowledgeBaseNodeData `json:"retrieveKnowledgeBaseNodeData"`
		StartNodeData                 *StartNodeData                 `json:"startNodeData"`
		EndNodeData                   *EndNodeData                   `json:"endNodeData"`
		CrawlerNodeData               *CrawlerNodeData               `json:"crawlerNodeData"`
		ConditionNodeData             *ConditionNodeData             `json:"conditionNodeData"`
	} `json:"data"`
}

// LLMNodeData LLM节点数据
type LLMNodeData struct {
	ModelName       string      `json:"model"`           // 模型名称
	ModelId         int64       `json:"modelId,string"`  // 模型ID
	Prompt          string      `json:"prompt"`          // 提示词
	InputVariables  []*Variable `json:"inputVariables"`  // 输入变量列表, key:变量名，value：变量来源{{nodeId.xxx}}或空(运行时输入)
	OutputFormat    string      `json:"outputFormat"`    // 输出格式 text,markdown,json
	OutputVariables []*Variable `json:"outputVariables"` // 输出变量名列表
	Temperature     float64     `json:"temperature"`     // 温度 0~2
	TopP            float64     `json:"topP"`            // TopP 0~1
}

// KnowledgeBaseWriteNodeData 写入知识库节点数据
type KnowledgeBaseWriteNodeData struct {
	KnowledgeBaseId int64  `json:"knowledgeBaseId"` // 知识库ID
	Content         string `json:"content"`         // 写入内容, {{nodeId.xxx}}表示写入某节点的变量列表的内容, 空表示运行时输入
}

// RetrieveKnowledgeBaseNodeData 检索知识库节点数据
type RetrieveKnowledgeBaseNodeData struct {
	KbId                int64        `json:"kbId,string"`         // 知识库ID
	SearchType          KbSearchType `json:"searchType"`          // 搜索类型
	Count               int          `json:"count"`               // 返回最大数量
	SimilarityThreshold float32      `json:"similarityThreshold"` // 相似度阈值
	OptimizeQuery       bool         `json:"optimizeQuery"`       // 是否优化用户输入
	InputVariables      []*Variable  `json:"inputVariables"`      // 输入变量，必填query（查询内容）
	OutputVariables     []*Variable  `json:"outputVariables"`     // 输出变量
}

type StartNodeData struct {
	InputVariables []*Variable `json:"inputVariables"` // 输入变量列表
}

type Variable struct {
	Type      string `json:"type"`      // 变量类型
	Name      string `json:"name"`      // 变量名
	Value     string `json:"value"`     // 变量值
	MustExist bool   `json:"mustExist"` // 是否必须存在
}

type EndNodeData struct {
	OutputVariables []*Variable `json:"outputVariables"` // 输出变量列表 key: 变量名，value: 变量来源 {{nodeId.xxx}}
}

type CrawlerNodeData struct {
	InputVariables  []*Variable `json:"inputVariables"`
	OutputVariables []*Variable `json:"outputVariables"`
}

type Condition struct {
	Value1 *Variable `json:"value1"`
	Value2 *Variable `json:"value2"`
	Op     string    `json:"op"`
}

type ConditionNodeBranch struct {
	Handle     string       `json:"handle"`
	Connector  string       `json:"connector"`
	Conditions []*Condition `json:"conditions"`
}

type ConditionNodeData struct {
	Branches []*ConditionNodeBranch `json:"branches"`
}

type ConditionNodeOutput struct {
	SuccessBranch string `json:"successBranch"`
}
