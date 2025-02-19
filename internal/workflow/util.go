package workflow

import (
	"github.com/StellrisJAY/workflow-ai/internal/model"
	"slices"
)

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
		if edge.Source == currNode.Id {
			childNodes = append(childNodes, nodes[edge.Target])
		}
	}
	return childNodes
}

func GetPrevNodes(definition *model.WorkflowDefinition, currNode *model.Node) []*model.Node {
	nodes := NodeSliceToMap(definition.Nodes)
	var prevNodes []*model.Node
	for _, edge := range definition.Edges {
		if edge.Target == currNode.Id {
			prevNodes = append(prevNodes, nodes[edge.Source])
		}
	}
	return prevNodes
}

func FindBranchNextNodes(definition *model.WorkflowDefinition, currNode *model.Node,
	branch *model.ConditionNodeBranch) []*model.Node {
	nodes := NodeSliceToMap(definition.Nodes)
	var nextNodes []*model.Node
	for _, edge := range definition.Edges {
		if edge.Source == currNode.Id && edge.SourceHandle == branch.Handle {
			nextNodes = append(nextNodes, nodes[edge.Target])
		}
	}
	return nextNodes
}

func FindNodeById(definition *model.WorkflowDefinition, id string) *model.Node {
	idx := slices.IndexFunc(definition.Nodes, func(x *model.Node) bool {
		return x.Id == id
	})
	if idx == -1 {
		return nil
	}
	return definition.Nodes[idx]
}

func FindNodeOutputVariable(node *model.Node, varName string) *model.Variable {
	var outputVars []*model.Variable
	switch model.NodeType(node.Type) {
	case model.NodeTypeStart:
		outputVars = node.Data.StartNodeData.InputVariables
	case model.NodeTypeLLM:
		outputVars = node.Data.LLMNodeData.OutputVariables
	case model.NodeTypeCrawler:
		outputVars = node.Data.CrawlerNodeData.OutputVariables
	case model.NodeTypeEnd:
		outputVars = node.Data.EndNodeData.OutputVariables
	}
	idx := slices.IndexFunc(outputVars, func(variable *model.Variable) bool {
		return variable.Name == varName
	})
	if idx == -1 {
		return nil
	}
	return outputVars[idx]
}

func GetPassedEdges(definition *model.WorkflowDefinition, nodes []*model.NodeStatusDTO, branches []string) []string {
	nodeMap := make(map[string]struct{})
	for _, node := range nodes {
		nodeMap[node.NodeId] = struct{}{}
	}
	branchMap := make(map[string]string)
	for _, branch := range branches {
		branchMap[branch] = branch
	}
	var passedEdges []string
	for _, edge := range definition.Edges {
		_, ok1 := nodeMap[edge.Target]
		_, ok2 := nodeMap[edge.Source]
		ok3 := true
		if edge.SourceHandle != "" {
			_, ok3 = branchMap[edge.SourceHandle]
		}
		if ok1 && ok2 && ok3 {
			passedEdges = append(passedEdges, edge.Id)
		}
	}
	return passedEdges
}
