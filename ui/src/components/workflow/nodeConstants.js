import {markRaw, ref} from "vue";
import LLMNode from "./node/llmNode.vue";
import StartNode from "./node/startNode.vue";
import EndNode from "./node/endNode.vue";
import conditionNode from "./node/conditionNode.vue";
import knowledgeRetrievalNode from "./node/knowledgeRetrivalNode.vue";
import knowledgeWriteNode from "./node/knowledgeWriteNode.vue";
import CustomEdge from "./customEdge.vue";
import CrawlerNode from "./node/crawlerNode.vue";
import templateAPI from "../../api/template.js";
import {randomUUID} from "../../util/uuid.js";
import webSearchNode from "./node/webSearchNode.vue";
import keywordExtractNode from "./node/keywordExtraction.vue";
import questionOptimizationNode from "./node/questionOptimizationNode.vue";
import imageUnderstandingNode from "./node/imageUnderstanding.vue";
import ocrNode from "./node/ocrNode.vue";

const nodeTypes = ref({
    llm: markRaw(LLMNode),
    start: markRaw(StartNode),
    end: markRaw(EndNode),
    condition: markRaw(conditionNode),
    knowledgeRetrieval: markRaw(knowledgeRetrievalNode),
    knowledgeWrite: markRaw(knowledgeWriteNode),
    crawler: markRaw(CrawlerNode),
    webSearch: markRaw(webSearchNode),
    keywordExtraction: markRaw(keywordExtractNode),
    questionOptimization: markRaw(questionOptimizationNode),
    imageUnderstanding: markRaw(imageUnderstandingNode),
    ocr: markRaw(ocrNode),
});

const edgeTypes = ref({
    custom: markRaw(CustomEdge),
});

const nodeTypeOptions = ref([
    { value: "llm", label: "大模型", description: "使用提示词和变量让大模型生成内容" },
    { value: "knowledgeRetrieval", label: "知识库检索", description: "从知识库检索数据" },
    // { value: "knowledgeWrite", label: "知识库写入", description: "将数据写入知识库" },
    { value: "condition", label: "条件", description: "条件分支" },
    { value: "crawler", label: "爬虫", description: "从给定地址爬取文本内容"},
    {value: "webSearch", label: "网页搜索", description: "使用搜索引擎搜索网页内容"},
    {value: "keywordExtraction", label: "关键词提取", description: "从问题提取关键词"},
    {value: "questionOptimization", label: "提问优化", description: "优化用户提出的问题，以适配不同场景"},
    {value: "imageUnderstanding", label: "图像理解", description: "根据提示词理解图片"},
    {value: "ocr", label: "文字提取", description: "从图片提取文字"}
]);

const llmOutputFormatOptions = [
    {label: "JSON", value: "JSON"},
    {label: "TEXT", value: "TEXT"},
];

function createBranch() {
    return {
        handle: randomUUID(),
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
    const prototype = JSON.parse(nodePrototypes.value[nodeType]);
    node.data = prototype.data;
}
const nodePrototypes = ref({});

async function loadNodePrototypes() {
    const prototypes = {};
    for (let i = 0; i < nodeTypeOptions.value.length; i++) {
        const nodeType = nodeTypeOptions.value[i];
        const resp = await templateAPI.getNodePrototype(nodeType.value);
        prototypes[nodeType.value] = resp.data;
    }
    nodePrototypes.value = prototypes;
}

export default {
    nodeTypes: nodeTypes,
    edgeTypes: edgeTypes,
    nodeTypeOptions: nodeTypeOptions,
    createNodeData: createNodeData,
    llmOutputFormatOptions: llmOutputFormatOptions,
    createBranch: createBranch,
    loadNodePrototypes: loadNodePrototypes,
    nodePrototypes: nodePrototypes,
    deleteNodeEvent: "delete-node-event",
    deleteEdgeEvent: "delete-edge-event",
};