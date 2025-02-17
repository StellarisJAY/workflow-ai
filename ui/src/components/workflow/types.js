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

function initLLMNodeData() {
    return {
        llmNodeData: {
            inputVariables: [{name:"input",type:"string",value:""}],
            outputVariables: [{name: "output", type: "string", value: ""}],
        }
    };
}

function initCrawlerNodeData() {
    return {
        crawlerNodeData: {
            inputVariables: [{name:"url",type:"string",value:"",mustExist: true}],
            outputVariables: [{name:"content",type:"string",value: "",mustExist: true}],
        }
    }
}

function createNodeData(node) {
    const nodeType = node.type;
    switch (nodeType) {
        case "llm": node.data = initLLMNodeData(); node.name="大模型"; break;
        case "knowledgeRetrieval": break;
        case "knowledgeWrite": break;
        case "condition": break;
        case "crawler": node.data = initCrawlerNodeData(); node.name="爬虫"; break;
    }
}

export default {
    nodeTypes: nodeTypes,
    edgeTypes: edgeTypes,
    nodeTypeOptions: nodeTypeOptions,
    initLLMNodeData: initLLMNodeData,
    createNodeData: createNodeData,
};