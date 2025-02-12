<script setup>
import { Handle, Position } from '@vue-flow/core';
import baseNode from './baseNode.vue';
import {FormItem, Input} from "ant-design-vue";

defineProps(['id', 'type', 'data', 'showControls']);
</script>

<template>
  <span>大模型调用</span>
	<Handle type="source" :position="Position.Left"></Handle>
	<Handle type="target" :position="Position.Right"></Handle>
	<baseNode :id="id" :show-controls="showControls" :type="type" :data="data"/>
  <FormItem label="模型">
    {{data.modelName}}
  </FormItem>
  <FormItem label="输入">{{data.inputVariables}}</FormItem>
  <FormItem label="输出">{{data.outputVariables}}</FormItem>
</template>

<style scoped></style>