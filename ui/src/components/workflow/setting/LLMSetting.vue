<script setup>
import {FormItem, Textarea, Select} from "ant-design-vue";
import VariableTable from "./VariableTable.vue";

defineProps(['node', 'llmList', 'refOptions', 'llmOptions']);
</script>

<template>
  <FormItem label="模型">
    <Select v-model:value="node.data['modelId']" :options="llmOptions"></Select>
  </FormItem>
  <FormItem label="系统提示词">
    <Textarea v-model:value="node.data['systemPrompt']" placeholder="输入系统提示词，使用{变量名}格式来嵌入变量"/>
  </FormItem>
  <FormItem label="用户提示词">
    <Textarea v-model:value="node.data['userPrompt']" placeholder="输入用户提示词，使用{变量名}格式来嵌入变量"/>
  </FormItem>
  <VariableTable :node-id="node.id"
                 :has-input="true"
                 :has-output="true"
                 :input-variables="node.data.inputVariables"
                 :output-variables="node.data.outputVariables"
                 :ref-options="refOptions"/>
</template>

<style scoped></style>