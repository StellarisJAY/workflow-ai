<script setup>
import {FormItem, Form, Input, Textarea} from "ant-design-vue";

defineProps(['node', 'llmList']);

</script>

<template>
  <FormItem label="模型">
    <Input v-model:value="node.data['modelName']"></Input>
  </FormItem>
  <FormItem label="系统提示词">
    <Textarea v-model:value="node.data['systemPrompt']" placeholder="输入系统提示词，使用{变量名}格式来嵌入变量"/>
  </FormItem>
  <FormItem label="用户提示词">
    <Textarea v-model:value="node.data['userPrompt']" placeholder="输入用户提示词，使用{变量名}格式来嵌入变量"/>
  </FormItem>
  <FormItem label="输入变量"></FormItem>
  <FormItem label="输出变量"></FormItem>
</template>

<style scoped></style>