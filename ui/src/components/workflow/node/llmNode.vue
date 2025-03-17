<script setup>
import {Handle, Position, useVueFlow} from '@vue-flow/core';
import {FormItem, Card} from "ant-design-vue";
import NodeExtra from "./nodeExtra.vue";
import NodeVariableDisplay from "./nodeVariableDisplay.vue";
import {onMounted} from "vue";
import NodeConstants from "../nodeConstants.js";
import NodeUtil from "../../../util/nodeUtil.js";

const props = defineProps(['id', 'type', 'data']);
const {findNode, getNodes, getEdges} = useVueFlow();
const node = findNode(props.id);

onMounted(()=>{
  addEventListener(NodeConstants.deleteNodeEvent, ev=>onNodeDelete(ev));
  addEventListener(NodeConstants.deleteEdgeEvent, ev=>onEdgeDelete(ev));
});

function onNodeDelete(ev) {
  const llmNodeData = props.data["llmNodeData"];
  NodeUtil.resetInputVariableRef(llmNodeData, id=>ev.detail.id === id);
}

function onEdgeDelete(ev) {
  // 获取所有前驱节点，通过判断变量引用节点是否在前驱节点列表判断是否需要重置引用
  const prevNodes = NodeUtil.getPrevNodes(props.id, getNodes.value, getEdges.value).map(node=>node.id);
  const llmNodeData = props.data["llmNodeData"];
  NodeUtil.resetInputVariableRef(llmNodeData, id=>!prevNodes.includes(id));
}
</script>

<template>
  <Card :title="node.data['name']" :hoverable="true" :body-style="{padding:'10px'}" :head-style="{padding:'10px'}">
    <Handle type="target" :position="Position.Left"></Handle>
    <Handle type="source" :position="Position.Right"></Handle>
    <template #extra>
      <node-extra :id="id" :type="type" :data="data" :status="node.status" :editable="true"/>
    </template>
    <FormItem label="模型">
      {{data['llmNodeData']['modelName']}}
    </FormItem>
    <node-variable-display :input-variables="data['input']"
                           :output-variables="data['output']"
                           :has-output="true"
                           :has-input="true"/>
  </Card>
</template>

<style scoped></style>