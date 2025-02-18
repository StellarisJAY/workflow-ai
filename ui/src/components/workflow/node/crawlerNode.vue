<script setup>
import {Handle, Position, useVueFlow} from '@vue-flow/core';
import {Card} from "ant-design-vue";
import NodeExtra from "./nodeExtra.vue";
import NodeVariableDisplay from "./nodeVariableDisplay.vue";

const props = defineProps(['id', 'type', 'data']);
const {findNode} = useVueFlow();
const node = findNode(props.id);
</script>

<template>
  <Handle type="source" :position="Position.Left"></Handle>
  <Handle type="target" :position="Position.Right"></Handle>
  <Card :title="node.data['name']" :hoverable="true">
    <template #extra>
      <node-extra :id="id" :type="type" :data="data" :status="node.status" :editable="true"/>
    </template>
    <node-variable-display :input-variables="data['crawlerNodeData'].inputVariables"
                           :output-variables="data['crawlerNodeData'].outputVariables"
                           :has-output="true"
                           :has-input="true"/>
  </Card>
</template>

<style scoped></style>