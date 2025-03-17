<script setup>
import {Handle, Position, useVueFlow} from '@vue-flow/core';
import {Card} from "ant-design-vue";
import NodeVariableDisplay from "./nodeVariableDisplay.vue";
import {onMounted} from "vue";
import NodeConstants from "../nodeConstants.js";
import NodeUtil from "../../../util/nodeUtil.js";
import NodeExtra from "./nodeExtra.vue";

const props = defineProps(['id', 'type', 'data']);
const {findNode, getEdges, getNodes} = useVueFlow();
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
      <node-extra :id="id" :type="type" :data="data" :status="node.status" :editable="true"/>
    </template>
    <Handle type="target" :position="Position.Left" :connectable="true"></Handle>
    <node-variable-display :has-input="true" :has-output="false"
                           :input-variables="data['input']"/>
  </Card>
</template>

<style scoped></style>