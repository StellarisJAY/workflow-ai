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
        }
        if (outputVariables) {
            let option = {
                label: node.data['name'],
                value: node.id,
                children: []
            };
            outputVariables.forEach(variable=>{
                option.children.push({label: variable.name, value: variable.name});
            });
            options.push(option);
        }
    });
    return options;
}

export default NodeUtil;