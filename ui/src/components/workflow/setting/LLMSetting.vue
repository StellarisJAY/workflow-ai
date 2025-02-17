<script setup>
import {FormItem, Textarea, Select} from "ant-design-vue";
import VariableTable from "./VariableTable.vue";

defineProps(['node', 'llmList', 'refOptions', 'llmOptions']);
</script>

<template>
  <FormItem label="模型">
    <Select v-model:value="node.data['llmNodeData']['modelId']" :options="llmOptions"></Select>
  </FormItem>
  <FormItem label="系统提示词">
    <Textarea v-model:value="node.data['llmNodeData']['systemPrompt']" placeholder="输入系统提示词，使用{变量名}格式来嵌入变量"/>
  </FormItem>
  <FormItem label="用户提示词">
    <Textarea v-model:value="node.data['llmNodeData']['userPrompt']" placeholder="输入用户提示词，使用{变量名}格式来嵌入变量"/>
  </FormItem>
  <VariableTable :node-id="node.id"
                 :has-input="true"
                 :has-output="true"
                 :input-variables="node.data['llmNodeData'].inputVariables"
                 :output-variables="node.data['llmNodeData'].outputVariables"
                 :ref-options="refOptions"
                 :input-editable="true"
                 :output-editable="true"/>
</template>

<style scoped></style>