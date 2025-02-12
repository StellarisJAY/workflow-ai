package workflow

import "github.com/StellrisJAY/workflow-ai/internal/model"

func NodeSliceToMap(nodes []*model.Node) map[string]*model.Node {
	nodeMap := make(map[string]*model.Node)
	for _, node := range nodes {
		nodeMap[node.Id] = node
	}
	return nodeMap
}

func GetNextNodes(definition *model.WorkflowDefinition, currNode *model.Node) []*model.Node {
	nodes := NodeSliceToMap(definition.Nodes)
	var childNodes []*model.Node
	for _, edge := range definition.Edges {
		if edge.Target == currNode.Id {
			childNodes = append(childNodes, nodes[edge.Source])
		}
	}
	return childNodes
}

func GetPrevNodes(definition *model.WorkflowDefinition, currNode *model.Node) []*model.Node {
	nodes := NodeSliceToMap(definition.Nodes)
	var prevNodes []*model.Node
	for _, edge := range definition.Edges {
		if edge.Source == currNode.Id {
			prevNodes = append(prevNodes, nodes[edge.Target])
		}
	}
	return prevNodes
}
