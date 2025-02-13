<script setup>
import { Handle, Position } from '@vue-flow/core';
import {FormItem, Card} from "ant-design-vue";
import NodeExtra from "./nodeExtra.vue";
import NodeVariableDisplay from "./nodeVariableDisplay.vue";

defineProps(['id', 'type', 'data', 'showControls']);
</script>

<template>
	<Handle type="source" :position="Position.Left"></Handle>
	<Handle type="target" :position="Position.Right"></Handle>
  <Card title="大模型" :hoverable="true">
    <template #extra>
      <node-extra :id="id" :show-controls="showControls" :type="type" :data="data" title="大模型"/>
    </template>
    <FormItem label="模型">
      {{data.modelName}}
    </FormItem>
    <node-variable-display :input-variables="data.inputVariables"
                           :output-variables="data.outputVariables"
                           :has-output="true"
                           :has-input="true"/>
  </Card>
</template>

<style scoped></style>