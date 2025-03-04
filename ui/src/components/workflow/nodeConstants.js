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
    // { value: "knowledgeWrite", label: "知识库写入", description: "将数据写入知识库" },
    { value: "condition", label: "条件", description: "条件分支" },
    { value: "crawler", label: "爬虫", description: "从给定地址爬取文本内容"},
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

const typeOptions = [
    {label: "string", value: "string"},
    {label: "文件", value: "file"},
    {label: "引用", value: "ref"},
];

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
    typeOptions: typeOptions,
    createBranch: createBranch,
    loadNodePrototypes: loadNodePrototypes,
    nodePrototypes: nodePrototypes,
    deleteNodeEvent: "delete-node-event",
    deleteEdgeEvent: "delete-edge-event",
};