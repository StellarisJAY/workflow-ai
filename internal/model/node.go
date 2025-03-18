package model

type NodeType string

const (
	NodeTypeStart                NodeType = "start"                // 开始节点
	NodeTypeLLM                  NodeType = "llm"                  // 大模型节点
	NodeTypeKnowledgeRetrieval   NodeType = "knowledgeRetrieval"   // 知识库检索节点
	NodeTypeKnowledgeWrite       NodeType = "knowledgeWrite"       // 写入知识库节点
	NodeTypeEnd                  NodeType = "end"                  // 结束节点
	NodeTypeCrawler              NodeType = "crawler"              // 爬虫节点
	NodeTypeCondition            NodeType = "condition"            // 条件判断节点
	NodeTypeKeywordExtraction    NodeType = "keywordExtraction"    // 关键词提取节点
	NodeTypeWebSearch            NodeType = "webSearch"            // 搜索引擎节点
	NodeTypeQuestionOptimization NodeType = "questionOptimization" // 问题优化节点
	NodeTypeImageUnderstanding   NodeType = "imageUnderstanding"   // 图像理解节点
	NodeTypeOCR                  NodeType = "ocr"                  // OCR文档识别节点
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
	Id       string   `json:"id"`   // 节点ID
	Type     NodeType `json:"type"` // 节点类型
	Position struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	} `json:"position"` // 节点位置
	Data NodeData `json:"data"`
}

type VarValueType string

const (
	VarValueTypeLiteral VarValueType = "literal" // 字面量
	VarValueTypeRef     VarValueType = "ref"     // 引用
)

type Value struct {
	Type       VarValueType `json:"type"`
	Content    string       `json:"content"`
	SourceNode string       `json:"sourceNode"`
	SourceName string       `json:"sourceName"`
}

type Input struct {
	Name     string       `json:"name"`
	Type     VariableType `json:"type"`
	Value    Value        `json:"value"`
	Required bool         `json:"required"`
	Fixed    bool         `json:"fixed"`
}

type Output struct {
	Name  string       `json:"name"`
	Type  VariableType `json:"type"`
	Value any          `json:"value"`
}

type NodeData struct {
	Name                 string         `json:"name"`
	AllowAddInputVar     bool           `json:"allowAddInputVar"`     // 是否允许添加输入变量
	AllowAddOutputVar    bool           `json:"allowAddOutputVar"`    // 是否允许添加输出变量
	DefaultAllowVarTypes []VariableType `json:"defaultAllowVarTypes"` // 默认允许的变量类型

	Input  []Input  `json:"input"`  // 节点输入变量列表
	Output []Output `json:"output"` // 节点输出变量列表

	LLMNodeData                   *LLMNodeData                   `json:"llmNodeData,omitempty"`                   // 大模型节点数据
	KnowledgeBaseWriteNodeData    *KnowledgeBaseWriteNodeData    `json:"knowledgeBaseWriteNodeData,omitempty"`    // 写入知识库节点数据
	RetrieveKnowledgeBaseNodeData *RetrieveKnowledgeBaseNodeData `json:"retrieveKnowledgeBaseNodeData,omitempty"` // 检索知识库节点数据
	StartNodeData                 *StartNodeData                 `json:"startNodeData,omitempty"`                 // 开始节点数据
	EndNodeData                   *EndNodeData                   `json:"endNodeData,omitempty"`                   // 结束节点数据
	CrawlerNodeData               *CrawlerNodeData               `json:"crawlerNodeData,omitempty"`               // 爬虫节点数据
	ConditionNodeData             *ConditionNodeData             `json:"conditionNodeData,omitempty"`             // 条件判断节点数据
	WebSearchNodeData             *WebSearchNodeData             `json:"webSearchNodeData,omitempty"`             // 搜索引擎节点数据
	KeywordExtractionNodeData     *KeywordExtractionNodeData     `json:"keywordExtractionNodeData,omitempty"`     // 关键词提取节点数据
	QuestionOptimizationNodeData  *QuestionOptimizationNodeData  `json:"questionOptimizationNodeData,omitempty"`  // 问题优化节点数据
	ImageUnderstandingNodeData    *ImageUnderstandingNodeData    `json:"imageUnderstandingNodeData,omitempty"`    // 图像理解节点数据
	OCRNodeData                   *OCRNodeData                   `json:"ocrNodeData,omitempty"`                   // OCR文档识别节点数据
}

// LLMNodeData LLM节点数据
type LLMNodeData struct {
	ModelName    string  `json:"modelName"`      // 模型名称
	ModelId      int64   `json:"modelId,string"` // 模型ID
	Prompt       string  `json:"prompt"`         // 提示词
	SystemPrompt string  `json:"systemPrompt"`   // 系统提示词
	UserPrompt   string  `json:"userPrompt"`     // 用户提示词
	OutputFormat string  `json:"outputFormat"`   // 输出格式 text,markdown,json
	Temperature  float64 `json:"temperature"`    // 温度 0~2
	TopP         float64 `json:"topP"`           // TopP 0~1
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
}

type StartNodeData struct {
}

type EndNodeData struct {
}

type CrawlerNodeData struct {
}

type Condition struct {
	Value1 *Input `json:"value1"`
	Value2 *Input `json:"value2"`
	Op     string `json:"op"`
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
	TopN int `json:"topN"` // 返回结果数量
}

type KeywordExtractionNodeData struct {
	ModelId   int64  `json:"modelId,string"`
	ModelName string `json:"modelName"`
	Count     int    `json:"count"`
}

type QuestionOptimizationNodeData struct {
	ModelId   int64  `json:"modelId,string"`
	ModelName string `json:"modelName"`
}

type ImageUnderstandingNodeData struct {
	ModelId      int64  `json:"modelId,string"`
	ModelName    string `json:"modelName"`
	Prompt       string `json:"prompt"`
	OutputFormat string `json:"outputFormat"`
}

type OCRNodeData struct {
	ModelId   int64  `json:"modelId,string"`
	ModelName string `json:"modelName"`
}

type MemoryNodeData struct {
}

// ConditionNodePrototype 条件判断节点原型
var ConditionNodePrototype = &Node{
	Type: NodeTypeCondition,
	Data: NodeData{
		Name:                 "条件",
		DefaultAllowVarTypes: []VariableType{VariableTypeString, VariableTypeNumber},
		ConditionNodeData: &ConditionNodeData{
			Branches: []*ConditionNodeBranch{
				{
					Handle:     "if",
					Connector:  "and",
					Conditions: []*Condition{},
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
	Type: NodeTypeLLM,
	Data: NodeData{
		Name:                 "大模型",
		DefaultAllowVarTypes: []VariableType{VariableTypeString, VariableTypeNumber},
		AllowAddInputVar:     true,
		AllowAddOutputVar:    true,
		Input:                []Input{},
		Output: []Output{
			{Name: "text", Type: VariableTypeString},
		},
		LLMNodeData: &LLMNodeData{
			SystemPrompt: "",
			UserPrompt:   "",
			OutputFormat: "JSON",
			Temperature:  0.5,
			TopP:         0.5,
			ModelName:    "",
			ModelId:      0,
		},
	},
}

var KbRetrievalNodePrototype = &Node{
	Type: NodeTypeKnowledgeRetrieval,
	Data: NodeData{
		Name:                 "知识库检索",
		DefaultAllowVarTypes: []VariableType{VariableTypeString, VariableTypeNumber},
		AllowAddInputVar:     false,
		AllowAddOutputVar:    false,
		Input: []Input{
			{Name: "query", Type: VariableTypeString, Required: true, Fixed: true},
		},
		Output: []Output{
			{Name: "total", Type: VariableTypeNumber},
			{Name: "documents", Type: VariableTypeStringArray},
		},
		RetrieveKnowledgeBaseNodeData: &RetrieveKnowledgeBaseNodeData{
			SearchType:          KbSearchTypeSimilarity,
			Count:               10,
			SimilarityThreshold: 0.8,
			OptimizeQuery:       false,
		},
	},
}

var CrawlerNodePrototype = &Node{
	Type: NodeTypeCrawler,
	Data: NodeData{
		Name:                 "HTTP请求",
		DefaultAllowVarTypes: []VariableType{VariableTypeString, VariableTypeNumber},
		AllowAddInputVar:     false,
		AllowAddOutputVar:    false,
		Input: []Input{
			{Name: "url", Type: VariableTypeString, Required: true, Fixed: true},
		},
		Output: []Output{
			{Name: "code", Type: VariableTypeNumber},
			{Name: "content-type", Type: VariableTypeString},
			{Name: "content", Type: VariableTypeString},
		},
		CrawlerNodeData: &CrawlerNodeData{},
	},
}

var KeywordExtractionNodePrototype = &Node{
	Type: NodeTypeKeywordExtraction,
	Data: NodeData{
		Name:                 "关键词提取",
		DefaultAllowVarTypes: []VariableType{VariableTypeString},
		AllowAddInputVar:     false,
		AllowAddOutputVar:    false,
		Input: []Input{
			{Name: "question", Type: VariableTypeString, Required: true, Fixed: true},
		},
		Output: []Output{
			{Name: "keywords", Type: VariableTypeStringArray},
		},
		KeywordExtractionNodeData: &KeywordExtractionNodeData{
			ModelId:   0,
			ModelName: "",
			Count:     3,
		},
	},
}

var QuestionOptimizationNodePrototype = &Node{
	Type: NodeTypeQuestionOptimization,
	Data: NodeData{
		Name:                 "问题优化",
		DefaultAllowVarTypes: []VariableType{VariableTypeString},
		AllowAddInputVar:     false,
		AllowAddOutputVar:    false,
		Input: []Input{
			{Name: "question", Type: VariableTypeString, Required: true, Fixed: true},
		},
		Output: []Output{
			{Name: "result", Type: VariableTypeString},
		},
		QuestionOptimizationNodeData: &QuestionOptimizationNodeData{
			ModelId:   0,
			ModelName: "",
		},
	},
}

var WebSearchNodePrototype = &Node{
	Type: NodeTypeWebSearch,
	Data: NodeData{
		Name:                 "网页搜索",
		DefaultAllowVarTypes: []VariableType{VariableTypeString},
		AllowAddInputVar:     false,
		AllowAddOutputVar:    false,
		Input: []Input{
			{Name: "query", Type: VariableTypeString, Required: true, Fixed: true},
		},
		Output: []Output{
			{Name: "urls", Type: VariableTypeStringArray},
			{Name: "contents", Type: VariableTypeStringArray},
		},
		WebSearchNodeData: &WebSearchNodeData{
			TopN: 10,
		},
	},
}

var ImageUnderstandingNodePrototype = &Node{
	Type: NodeTypeImageUnderstanding,
	Data: NodeData{
		Name:                 "图像理解",
		DefaultAllowVarTypes: []VariableType{VariableTypeString},
		AllowAddInputVar:     false,
		AllowAddOutputVar:    true,
		Input: []Input{
			{Name: "image", Type: VariableTypeImageFile, Required: true, Fixed: true},
		},
		Output: []Output{
			{Name: "text", Type: VariableTypeString},
		},
		ImageUnderstandingNodeData: &ImageUnderstandingNodeData{
			ModelId:      0,
			ModelName:    "",
			Prompt:       "",
			OutputFormat: "JSON",
		},
	},
}

var OCRNodePrototype = &Node{
	Type: NodeTypeOCR,
	Data: NodeData{
		Name:              "图片文字提取",
		AllowAddInputVar:  false,
		AllowAddOutputVar: false,
		Input: []Input{
			{Name: "image", Type: VariableTypeImageFile, Required: true, Fixed: true},
		},
		Output: []Output{
			{Name: "text", Type: VariableTypeString},
		},
		OCRNodeData: &OCRNodeData{},
	},
}

var EndNodePrototype = &Node{
	Type: NodeTypeEnd,
	Data: NodeData{
		Name:              "结束",
		AllowAddInputVar:  true,
		AllowAddOutputVar: false,
		Input:             []Input{},
		Output:            []Output{},
		EndNodeData:       &EndNodeData{},
	},
}

var SystemVariablePrototype = []Input{
	{Name: "sys.query", Type: VariableTypeString, Fixed: true, Value: Value{Type: VarValueTypeLiteral, Content: ""}},
	{Name: "sys.workflow_id", Type: VariableTypeString, Fixed: true, Value: Value{Type: VarValueTypeLiteral, Content: ""}},
}
