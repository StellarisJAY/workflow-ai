package model

type NodeType string

const (
	NodeTypeStart              NodeType = "start"              // 开始节点
	NodeTypeLLM                NodeType = "llm"                // 大模型节点
	NodeTypeKnowledgeRetrieval NodeType = "knowledgeRetrieval" // 知识库检索节点
	NodeTypeKnowledgeWrite     NodeType = "knowledgeWrite"     // 写入知识库节点
	NodeTypeEnd                NodeType = "end"                // 结束节点
	NodeTypeCrawler            NodeType = "crawler"            // 爬虫节点
	NodeTypeCondition          NodeType = "condition"          // 条件判断节点
	NodeTypeKeywordExtraction  NodeType = "keywordExtraction"  // 关键词提取节点
	NodeTypeWebSearch          NodeType = "webSearch"          // 搜索引擎节点
)

type VariableType string

const (
	VariableTypeNumber      VariableType = "number"     // 数值类型
	VariableTypeString      VariableType = "string"     // 字符串类型
	VariableTypeStringArray VariableType = "array_str"  // 字符串数组类型
	VariableTypeNumberArray VariableType = "array_num"  // 数值数组类型
	VariableTypeTextFile    VariableType = "text_file"  // 文本文件类型
	VariableTypeImageFile   VariableType = "image_file" // 图片文件类型
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
	AllowAddInputVar              bool                           `json:"allowAddInputVar"`              // 是否允许添加输入变量
	AllowAddOutputVar             bool                           `json:"allowAddOutputVar"`             // 是否允许添加输出变量
	DefaultAllowVarTypes          []VariableType                 `json:"defaultAllowVarTypes"`          // 默认允许的变量类型
	LLMNodeData                   *LLMNodeData                   `json:"llmNodeData"`                   // 大模型节点数据
	KnowledgeBaseWriteNodeData    *KnowledgeBaseWriteNodeData    `json:"knowledgeBaseWriteNodeData"`    // 写入知识库节点数据
	RetrieveKnowledgeBaseNodeData *RetrieveKnowledgeBaseNodeData `json:"retrieveKnowledgeBaseNodeData"` // 检索知识库节点数据
	StartNodeData                 *StartNodeData                 `json:"startNodeData"`                 // 开始节点数据
	EndNodeData                   *EndNodeData                   `json:"endNodeData"`                   // 结束节点数据
	CrawlerNodeData               *CrawlerNodeData               `json:"crawlerNodeData"`               // 爬虫节点数据
	ConditionNodeData             *ConditionNodeData             `json:"conditionNodeData"`             // 条件判断节点数据
	WebSearchNodeData             *WebSearchNodeData             `json:"webSearchNodeData"`             // 搜索引擎节点数据
}

// LLMNodeData LLM节点数据
type LLMNodeData struct {
	ModelName       string      `json:"modelName"`       // 模型名称
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
	KbId       int64    `json:"kbId,string"` // 知识库ID
	ChunkSize  int      `json:"chunkSize"`   // 分片大小
	Separators []string `json:"separators"`  // 分隔符
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
	Value        string         `json:"value"`        // 变量值, 文件类型的value为文件id
	Ref          string         `json:"ref"`          // 引用变量名，引用节点实例ID/变量名，只能
	AllowedTypes []VariableType `json:"allowedTypes"` // 允许的变量类型
	AllowRef     bool           `json:"allowRef"`     // 是否允许引用
	IsRef        bool           `json:"isRef"`        // 是否引用
	Required     bool           `json:"required"`     // 是否必填, 必填后不可删除
	Fixed        bool           `json:"fixed"`        // 是否固定, 固定后不可修改
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
	Handle     string       `json:"handle"`    // 分支唯一ID
	Connector  string       `json:"connector"` // 条件连接 and or
	Conditions []*Condition `json:"conditions"`
}

type ConditionNodeData struct {
	Branches []*ConditionNodeBranch `json:"branches"`
}

type ConditionNodeOutput struct {
	SuccessBranch string `json:"successBranch"`
}

type WebSearchType string

type WebSearchNodeData struct {
	TopN            int         `json:"topN"` // 返回结果数量
	InputVariables  []*Variable `json:"inputVariables"`
	OutputVariables []*Variable `json:"outputVariables"`
}

// ConditionNodePrototype 条件判断节点原型
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

// LLMNodePrototype 大模型节点原型
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
			ModelName:       "",
			ModelId:         0,
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
				{Name: "url", Type: VariableTypeString, Required: true, AllowRef: true, AllowedTypes: []VariableType{VariableTypeString}}, // 网页url
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

var KeywordExtractionNodePrototype = &Node{
	Type: string(NodeTypeKeywordExtraction),
	Data: NodeData{
		Name:                 "关键词提取",
		DefaultAllowVarTypes: []VariableType{VariableTypeString},
		AllowAddInputVar:     false,
		AllowAddOutputVar:    false,
		CrawlerNodeData: &CrawlerNodeData{
			InputVariables: []*Variable{
				{Name: "question", Type: VariableTypeString, Required: true, AllowRef: true}, // 问题
			},
			OutputVariables: []*Variable{
				{Name: "total", Type: VariableTypeNumber, Required: true, Fixed: true},         // 关键词数量
				{Name: "keywords", Type: VariableTypeStringArray, Required: true, Fixed: true}, // 关键词列表
			},
		},
	},
}

var WebSearchNodePrototype = &Node{
	Type: string(NodeTypeWebSearch),
	Data: NodeData{
		Name:                 "网页搜索",
		DefaultAllowVarTypes: []VariableType{VariableTypeString},
		AllowAddInputVar:     false,
		AllowAddOutputVar:    false,
		WebSearchNodeData: &WebSearchNodeData{
			TopN: 10,
			InputVariables: []*Variable{
				{Name: "query", Type: VariableTypeString, Required: true, Fixed: true, AllowRef: true},
			},
			OutputVariables: []*Variable{
				{Name: "total", Type: VariableTypeNumber, Required: true, Fixed: true},         // 搜索到的网页总数
				{Name: "urls", Type: VariableTypeStringArray, Required: true, Fixed: true},     // 搜索到的网页url列表
				{Name: "contents", Type: VariableTypeStringArray, Required: true, Fixed: true}, // 搜索到网页内容列表
			},
		},
	},
}
