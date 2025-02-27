<script setup>
import {Handle, Position, useVueFlow} from '@vue-flow/core';
import {Card} from "ant-design-vue";
import NodeVariableDisplay from "./nodeVariableDisplay.vue";
import NodeStatusTag from "./nodeStatusTag.vue";
import {onMounted} from "vue";
import NodeConstants from "../nodeConstants.js";
import NodeUtil from "../../../util/nodeUtil.js";

const props = defineProps(['id', 'type', 'data']);
const {findNode, nodesDraggable, nodesConnectable, getEdges, getNodes} = useVueFlow();
const node = findNode(props.id);

onMounted(()=>{
  addEventListener(NodeConstants.deleteNodeEvent, ev=>onNodeDelete(ev));
  addEventListener(NodeConstants.deleteEdgeEvent, ev=>onEdgeDelete(ev));
});

function onNodeDelete(ev) {
  const crawlerNodeData = props.data["endNodeData"];
  NodeUtil.resetOutputVariableRef(crawlerNodeData, id=>ev.detail.id === id);
}

function onEdgeDelete(ev) {
  // 获取所有前驱节点，通过判断变量引用节点是否在前驱节点列表判断是否需要重置引用
  const prevNodes = NodeUtil.getPrevNodes(props.id, getNodes.value, getEdges.value).map(node=>node.id);
  const crawlerNodeData = props.data["endNodeData"];
  NodeUtil.resetOutputVariableRef(crawlerNodeData, id=>!prevNodes.includes(id));
}

</script>

<template>
  <Card title="结束" :hoverable="true" :body-style="{padding:'10px'}" :head-style="{padding:'10px'}">
    <template #extra>
      <NodeStatusTag v-if="!nodesDraggable && !nodesConnectable" :status="node.status"/>
    </template>
    <Handle type="target" :position="Position.Left" :connectable="true"></Handle>
    <node-variable-display :has-input="false" :has-output="true"
                           :output-variables="data['endNodeData'].outputVariables"/>
  </Card>
</template>

<style scoped></style>