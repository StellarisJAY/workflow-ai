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
	VariableTypeNumber      VariableType = "number"
	VariableTypeString      VariableType = "string"
	VariableTypeStringArray VariableType = "array_str"
	VariableTypeNumberArray VariableType = "array_num"
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
	Data NodeData `json:"data"`
}

type NodeData struct {
	Name                          string                         `json:"name"`
	AllowAddInputVar              bool                           `json:"allowAddInputVar"`
	AllowAddOutputVar             bool                           `json:"allowAddOutputVar"`
	DefaultAllowVarTypes          []VariableType                 `json:"defaultAllowVarTypes"`
	LLMNodeData                   *LLMNodeData                   `json:"llmNodeData"`
	KnowledgeBaseWriteNodeData    *KnowledgeBaseWriteNodeData    `json:"knowledgeBaseWriteNodeData"`
	RetrieveKnowledgeBaseNodeData *RetrieveKnowledgeBaseNodeData `json:"retrieveKnowledgeBaseNodeData"`
	StartNodeData                 *StartNodeData                 `json:"startNodeData"`
	EndNodeData                   *EndNodeData                   `json:"endNodeData"`
	CrawlerNodeData               *CrawlerNodeData               `json:"crawlerNodeData"`
	ConditionNodeData             *ConditionNodeData             `json:"conditionNodeData"`
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
	Type         VariableType   `json:"type"`         // 变量类型
	Name         string         `json:"name"`         // 变量名
	Value        string         `json:"value"`        // 变量值
	Ref          string         `json:"ref"`          // 引用变量名，引用节点实例ID/变量名，只能
	AllowedTypes []VariableType `json:"allowedTypes"` // 允许的变量类型
	AllowRef     bool           `json:"allowRef"`
	IsRef        bool           `json:"isRef"`
	Required     bool           `json:"required"` // 是否必填, 必填后不可删除
	Fixed        bool           `json:"fixed"`    // 是否固定, 固定后不可修改
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

var ConditionNodePrototype = &Node{
	Type: string(NodeTypeCondition),
	Data: NodeData{
		Name:                 "条件",
		DefaultAllowVarTypes: []VariableType{VariableTypeString, VariableTypeNumber},
		ConditionNodeData: &ConditionNodeData{
			Branches: []*ConditionNodeBranch{
				{
					Handle:    "if",
					Connector: "and",
					Conditions: []*Condition{
						{
							Value1: &Variable{Value: "0", Type: VariableTypeNumber},
							Value2: &Variable{Value: "0", Type: VariableTypeNumber},
							Op:     "==",
						},
					},
				},
				{
					Handle: "else",
				},
			},
		},
	},
}

var LLMNodePrototype = &Node{
	Type: string(NodeTypeLLM),
	Data: NodeData{
		Name:                 "大模型",
		DefaultAllowVarTypes: []VariableType{VariableTypeString, VariableTypeNumber},
		AllowAddInputVar:     true,
		AllowAddOutputVar:    true,
		LLMNodeData: &LLMNodeData{
			Prompt: "",
			InputVariables: []*Variable{
				{Name: "input", Value: "", Type: VariableTypeString, AllowRef: true},
			},
			OutputFormat:    "JSON",
			OutputVariables: []*Variable{},
			Temperature:     0.5,
			TopP:            0.5,
		},
	},
}

var KbRetrievalNodePrototype = &Node{
	Type: string(NodeTypeKnowledgeRetrieval),
	Data: NodeData{
		Name:                 "知识库检索",
		DefaultAllowVarTypes: []VariableType{VariableTypeString, VariableTypeNumber},
		AllowAddInputVar:     false,
		AllowAddOutputVar:    false,
		RetrieveKnowledgeBaseNodeData: &RetrieveKnowledgeBaseNodeData{
			SearchType:          KbSearchTypeSimilarity,
			Count:               10,
			SimilarityThreshold: 0.8,
			OptimizeQuery:       false,
			InputVariables: []*Variable{
				{Name: "query", Type: VariableTypeString, Required: true, Fixed: true, AllowRef: true}, // 检索内容
			},
			OutputVariables: []*Variable{
				{Name: "total", Type: VariableTypeNumber, Required: true, Fixed: true},          // 检索到的文档总数
				{Name: "documents", Type: VariableTypeStringArray, Required: true, Fixed: true}, // 文档列表
			},
		},
	},
}

var CrawlerNodePrototype = &Node{
	Type: string(NodeTypeCrawler),
	Data: NodeData{
		Name:                 "HTTP请求",
		DefaultAllowVarTypes: []VariableType{VariableTypeString, VariableTypeNumber},
		AllowAddInputVar:     false,
		AllowAddOutputVar:    false,
		CrawlerNodeData: &CrawlerNodeData{
			InputVariables: []*Variable{
				{Name: "url", Type: VariableTypeString, Required: true, AllowRef: true}, // 网页url
			},
			OutputVariables: []*Variable{
				{Name: "code", Type: VariableTypeNumber, Required: true, Fixed: true},        // HTTP状态码
				{Name: "message", Type: VariableTypeString, Required: true, Fixed: true},     // HTTP状态码描述
				{Name: "data", Type: VariableTypeString, Required: true, Fixed: true},        // 网页内容
				{Name: "contentType", Type: VariableTypeString, Required: true, Fixed: true}, // 网页内容类型

			},
		},
	},
}
