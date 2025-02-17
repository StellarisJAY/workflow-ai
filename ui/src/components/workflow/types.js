import {markRaw, ref} from "vue";
import LLMNode from "./node/llmNode.vue";
import StartNode from "./node/startNode.vue";
import EndNode from "./node/endNode.vue";
import conditionNode from "./node/conditionNode.vue";
import knowledgeRetrievalNode from "./node/knowledgeRetrivalNode.vue";
import knowledgeWriteNode from "./node/knowledgeWriteNode.vue";
import CustomEdge from "./customEdge.vue";

const nodeTypes = ref({
    llm: markRaw(LLMNode),
    start: markRaw(StartNode),
    end: markRaw(EndNode),
    condition: markRaw(conditionNode),
    knowledgeRetrieval: markRaw(knowledgeRetrievalNode),
    knowledgeWrite: markRaw(knowledgeWriteNode)
});

const edgeTypes = ref({
    custom: markRaw(CustomEdge),
});

export default {
    nodeTypes: nodeTypes,
    edgeTypes: edgeTypes,
};