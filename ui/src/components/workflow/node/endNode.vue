<script setup>
import {Handle, Position} from '@vue-flow/core';
import {Card} from "ant-design-vue";
import NodeVariableDisplay from "./nodeVariableDisplay.vue";

const props = defineProps(['id', 'type', 'data']);
</script>

<template>
  <Card title="结束" :hoverable="true">
    <Handle type="source" :position="Position.Left" :connectable="true"></Handle>
    <node-variable-display :has-input="false" :has-output="true"
                           :output-variables="data['endNodeData'].outputVariables"/>
  </Card>
</template>

<style scoped></style>