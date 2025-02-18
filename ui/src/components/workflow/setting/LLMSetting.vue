<script setup>
import {FormItem, Textarea, Select, Slider, Input} from "ant-design-vue";
import VariableTable from "./VariableTable.vue";
import types from "../types.js";
import CommonSetting from "./CommonSetting.vue";
const props = defineProps(['node', 'llmList', 'refOptions', 'llmOptions']);

function onModelChange(modelId) {
  const llm = props.llmList.find(llm=>llm.id === modelId);
  if (llm) {
    props.node.data['llmNodeData']['modelName'] = llm.name;
  }
}
</script>

<template>
  <CommonSetting :node="node"/>
  <FormItem label="模型">
    <Select v-model:value="node.data['llmNodeData']['modelId']" :options="llmOptions" @change="onModelChange"></Select>
  </FormItem>
  <FormItem label="提示词">
    <Textarea class="prompt-textarea" v-model:value="node.data['llmNodeData']['prompt']" placeholder="输入提示词，使用{\{.变量名}}格式来嵌入变量"/>
  </FormItem>
  <FormItem label="温度">
    <Slider v-model:value="node.data['llmNodeData']['temperature']" :min="0" :max="2" :step="0.1"></Slider>
  </FormItem>
  <FormItem label="TopP">
    <Slider v-model:value="node.data['llmNodeData']['topP']" :min="0" :max="1" :step="0.1"></Slider>
  </FormItem>
  <FormItem label="输出格式">
    <Select v-model:value="node.data['llmNodeData']['outputFormat']" :options="types.llmOutputFormatOptions"/>
  </FormItem>
  <VariableTable :node-id="node.id"
                 :has-input="true"
                 :has-output="true"
                 :input-variables="node.data['llmNodeData'].inputVariables"
                 :output-variables="node.data['llmNodeData'].outputVariables"
                 :ref-options="refOptions"
                 :input-editable="true"
                 :output-editable="true" :allow-add-del-output="true" :allow-add-del-input="true"/>
</template>

<style scoped>
.prompt-textarea {
  height: 200px;
}
</style>