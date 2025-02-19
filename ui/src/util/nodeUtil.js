const NodeUtil = {
    getPrevNodesRecursive: function(nodeId, nodes, edges, prevNodes) {
        const prevEdges = edges.filter(edge=>edge.target === nodeId);
        if (prevEdges.length === 0) {
            return;
        }
        const targetNodeIds = prevEdges.flatMap(edge=> {return edge.source;});
        const parentNodes = nodes.filter(node=>targetNodeIds.includes(node.id));
        parentNodes.forEach(node=>prevNodes.push(node));
        targetNodeIds.forEach(id=>this.getPrevNodesRecursive(id, nodes, edges, prevNodes));
    },
    getPrevNodes: function(nodeId, nodes, edges) {
        let prevNodes = [];
        this.getPrevNodesRecursive(nodeId, nodes, edges, prevNodes);
        return Array.from(new Set(prevNodes));
    }
}

export default NodeUtil;