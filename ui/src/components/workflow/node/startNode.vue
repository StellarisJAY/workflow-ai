<script setup>
import {Handle, Position, useVueFlow} from '@vue-flow/core';
import {Card} from "ant-design-vue";
import NodeVariableDisplay from "./nodeVariableDisplay.vue";
import NodeExtra from "./nodeExtra.vue";

const props = defineProps(['id', 'type', 'data']);
const {findNode} = useVueFlow();
const node = findNode(props.id);
</script>

<template>
  <Card title="开始" :hoverable="true" body-style="padding:10px" head-style="padding:10px;">
    <Handle type="target" :position="Position.Right" :connectable="true"/>
    <template #extra>
      <node-extra :editable="false" :status="node.status"></node-extra>
    </template>
    <node-variable-display :has-input="true"
                           :has-output="false"
                           :input-variables="data['startNodeData'].inputVariables"></node-variable-display>
  </Card>
</template>

<style scoped></style>