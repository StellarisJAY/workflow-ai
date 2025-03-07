import {useVueFlow} from "@vue-flow/core";

const NodeUtil = {
    /**
     * 获取所有父节点
     * @param nodeId 节点id
     * @param nodes 所有节点列表
     * @param edges 所有连线列表
     * @returns {*[]} 父节点列表
     */
    getPrevNodes: function(nodeId, nodes, edges) {
        let prevNodes = [];
        getPrevNodesRecursive(nodeId, nodes, edges, prevNodes);
        return Array.from(new Set(prevNodes));
    },
    getPrevNodesOutputs: getPrevNodesOutputs,
    resetInputVariableRef: resetInputVariableRef,
    resetOutputVariableRef: resetOutputVariableRef,
}

function getPrevNodesRecursive(nodeId, nodes, edges, prevNodes) {
    const prevEdges = edges.filter(edge=>edge.target === nodeId);
    if (prevEdges.length === 0) {
        return;
    }
    const targetNodeIds = prevEdges.flatMap(edge=> {return edge.source;});
    const parentNodes = nodes.filter(node=>targetNodeIds.includes(node.id));
    parentNodes.forEach(node=>prevNodes.push(node));
    targetNodeIds.forEach(id=>getPrevNodesRecursive(id, nodes, edges, prevNodes));
}

/**
 * 获取所有父节点的输出变量
 * @param currNodeId 当前节点id
 * @returns {*[]} 父节点输出变量列表
 */
function getPrevNodesOutputs(currNodeId) {
    const {getNodes, getEdges} = useVueFlow();
    const prevNodes = NodeUtil.getPrevNodes(currNodeId, getNodes.value, getEdges.value);
    let options = [];
    prevNodes.forEach(node=>{
        if (!node.data) return;
        let outputVariables;
        switch (node.type) {
            case "llm": outputVariables = node.data['llmNodeData'].outputVariables; break;
            case "start": outputVariables = node.data['startNodeData'].inputVariables; break;
            case "end": outputVariables = node.data['endNodeData'].outputVariables; break;
            case "crawler": outputVariables = node.data['crawlerNodeData'].outputVariables; break;
            case "knowledgeRetrieval": outputVariables = node.data['retrieveKnowledgeBaseNodeData'].outputVariables; break;
            case "webSearch": outputVariables = node.data['webSearchNodeData'].outputVariables; break;
            case "keywordExtraction": outputVariables = node.data['keywordExtractionNodeData'].outputVariables; break;
            case "questionOptimization": outputVariables = node.data['questionOptimizationNodeData'].outputVariables; break;
            case "imageUnderstanding": outputVariables = node.data['imageUnderstandingNodeData'].outputVariables; break;
            case "ocr": outputVariables = node.data['ocrNodeData'].outputVariables; break;
        }
        if (outputVariables) {
            let option = {
                label: node.data['name'],
                value: node.id,
                children: [],
            };
            outputVariables.forEach(variable=>{
                option.children.push({label: variable.name, value: variable.name, type: variable.type});
            });
            options.push(option);
        }
    });
    return options;
}

/**
 * 删除节点或连线导致引用失效，重置节点输入变量列表的引用
 * @param nodeData 节点数据
 * @param needReset 通过变量引用的节点id判断是否需要重置的函数
 */
function resetInputVariableRef(nodeData, needReset) {
    nodeData.inputVariables.forEach(variable=>variableReset(variable, needReset));
}
/**
 * 删除节点或连线导致引用失效，重置节点输出变量列表的引用
 * @param nodeData 节点数据
 * @param needReset 通过变量引用的节点id判断是否需要重置的函数
 */
function resetOutputVariableRef(nodeData, needReset) {
    nodeData.outputVariables.forEach(variable=>variableReset(variable, needReset));
}

function variableReset(variable, needReset) {
    if (variable.isRef) {
        const ref = variable.ref.split(".");
        if (ref.length === 2 && needReset(ref[0])) {
            variable.isRef = true;
            variable.ref = "";
            variable.refOption = [];
        }
    }
}

export default NodeUtil;