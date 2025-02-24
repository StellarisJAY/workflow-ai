<script setup>
import {Handle, Position, useVueFlow} from '@vue-flow/core';
import NodeVariableDisplay from "./nodeVariableDisplay.vue";
import {Card} from "ant-design-vue";
import NodeExtra from "./nodeExtra.vue";

const props = defineProps(['id', 'type', 'data']);
const {findNode} = useVueFlow();
const node = findNode(props.id);
</script>

<template>
  <Card :title="node.data['name']" :hoverable="true" :body-style="{padding:'10px'}" :head-style="{padding:'10px'}">
    <Handle type="target" :position="Position.Left"></Handle>
    <Handle type="source" :position="Position.Right"></Handle>
    <template #extra>
      <node-extra :id="id" :type="type" :data="data" :status="node.status" :editable="true"/>
    </template>
    <node-variable-display :input-variables="data['knowledgeRetrievalNodeData'].inputVariables"
                           :output-variables="data['knowledgeRetrievalNodeData'].outputVariables"
                           :has-output="true"
                           :has-input="true"/>
  </Card>
</template>

<style scoped></style>