<script setup>
import {Handle, Position, useVueFlow} from '@vue-flow/core';
import {Card} from "ant-design-vue";
import NodeVariableDisplay from "./nodeVariableDisplay.vue";
import NodeStatusTag from "./nodeStatusTag.vue";

const props = defineProps(['id', 'type', 'data']);
const {findNode, nodesDraggable, nodesConnectable} = useVueFlow();
const node = findNode(props.id);
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