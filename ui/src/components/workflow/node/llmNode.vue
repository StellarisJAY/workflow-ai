<script setup>
import {Handle, Position, useVueFlow} from '@vue-flow/core';
import {FormItem, Card} from "ant-design-vue";
import NodeExtra from "./nodeExtra.vue";
import NodeVariableDisplay from "./nodeVariableDisplay.vue";

const props = defineProps(['id', 'type', 'data']);
const {findNode} = useVueFlow();
const node = findNode(props.id);
</script>

<template>
  <Card :title="node.data['name']" :hoverable="true">
    <Handle type="source" :position="Position.Left"></Handle>
    <Handle type="target" :position="Position.Right"></Handle>
    <template #extra>
      <node-extra :id="id" :type="type" :data="data" :status="node.status" :editable="true"/>
    </template>
    <FormItem label="模型">
      {{data['llmNodeData']['modelName']}}
    </FormItem>
    <node-variable-display :input-variables="data['llmNodeData'].inputVariables"
                           :output-variables="data['llmNodeData'].outputVariables"
                           :has-output="true"
                           :has-input="true"/>
  </Card>
</template>

<style scoped></style>