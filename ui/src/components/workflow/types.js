import {markRaw, ref} from "vue";
import LLMNode from "./node/llmNode.vue";
import StartNode from "./node/startNode.vue";
import EndNode from "./node/endNode.vue";
import conditionNode from "./node/conditionNode.vue";
import knowledgeRetrievalNode from "./node/knowledgeRetrivalNode.vue";
import knowledgeWriteNode from "./node/knowledgeWriteNode.vue";
import CustomEdge from "./customEdge.vue";
import CrawlerNode from "./node/crawlerNode.vue";

const nodeTypes = ref({
    llm: markRaw(LLMNode),
    start: markRaw(StartNode),
    end: markRaw(EndNode),
    condition: markRaw(conditionNode),
    knowledgeRetrieval: markRaw(knowledgeRetrievalNode),
    knowledgeWrite: markRaw(knowledgeWriteNode),
    crawler: markRaw(CrawlerNode),
});

const edgeTypes = ref({
    custom: markRaw(CustomEdge),
});

const nodeTypeOptions = ref([
    { value: "llm", label: "大模型", description: "使用提示词和变量让大模型生成内容" },
    { value: "knowledgeRetrieval", label: "知识库检索", description: "从知识库检索数据" },
    { value: "knowledgeWrite", label: "知识库写入", description: "将数据写入知识库" },
    { value: "condition", label: "条件", description: "条件分支" },
    { value: "crawler", label: "爬虫", description: "从给定地址爬取文本内容"},
]);

const llmOutputFormatOptions = [
    {label: "JSON", value: "JSON"},
    {label: "TEXT", value: "TEXT"},
];

function initLLMNodeData() {
    return {
        name: "大模型",
        llmNodeData: {
            inputVariables: [{name:"input",type:"string",value:""}],
            outputVariables: [{name: "output", type: "string", value: ""}],
            temperature: 0.5,
            topP: 0.5,
            outputFormat: "TEXT",
        }
    };
}

function initCrawlerNodeData() {
    return {
        name: "爬虫",
        crawlerNodeData: {
            inputVariables: [{name:"url",type:"string",value:"",mustExist: true}],
            outputVariables: [
                {name:"data",type:"string",value: "",mustExist: true},
                {name:"message",type:"string",value: "",mustExist: true},
                {name:"code", type:"string",value:"",mustExist: true},
                {name:"contentType", type:"string",value:"",mustExist: true},
            ],
        }
    }
}

function initConditionNodeData() {
    return {
        name: "条件判断",
        conditionNodeData: {
            branches: [
                {
                    handle: "1",
                    connector: "and",
                    conditions: [
                        {
                            value1: {type: "string", value: "value", refOption: []},
                            op: "==",
                            value2: {type: "string", value: "value", refOption: []},
                        }
                    ]
                },
                {
                    handle: "else"
                },
            ]
        }
    }
}

function createBranch() {
    return {
        handle: crypto.randomUUID(),
        connector: "and",
        conditions: [
            {
                value1: {type: "string", value: "value", refOption: []},
                op: "==",
                value2: {type: "string", value: "value", refOption: []},
            }
        ]
    }
}

function createNodeData(node) {
    const nodeType = node.type;
    switch (nodeType) {
        case "llm": node.data = initLLMNodeData(); break;
        case "knowledgeRetrieval": break;
        case "knowledgeWrite": break;
        case "condition": node.data = initConditionNodeData(); break;
        case "crawler": node.data = initCrawlerNodeData(); break;
    }
}

const typeOptions = [
    {label: "string", value: "string"},
    {label: "文件", value: "file"},
    {label: "引用", value: "ref"},
];

export default {
    nodeTypes: nodeTypes,
    edgeTypes: edgeTypes,
    nodeTypeOptions: nodeTypeOptions,
    initLLMNodeData: initLLMNodeData,
    createNodeData: createNodeData,
    llmOutputFormatOptions: llmOutputFormatOptions,
    typeOptions: typeOptions,
    createBranch: createBranch,
};