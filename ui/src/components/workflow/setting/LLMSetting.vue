<script setup>
import {FormItem, Textarea, Select, Slider} from "ant-design-vue";
import VariableTable from "./VariableTable.vue";
import nodeConstants from "../nodeConstants.js";
import CommonSetting from "./CommonSetting.vue";
import {onMounted, ref} from "vue";
import llmAPI from "../../../api/llm.js";
const props = defineProps(['node']);

const llmList = ref([]);
const llmOptions = ref([]);

onMounted(()=>{
  llmAPI.listModels({paged: false}).then(resp=>{
    llmList.value = resp.data;
    const options = [];
    llmList.value.forEach(item=>{
      options.push({label: item.name, value: item.id});
    });
    llmOptions.value = options;
  });
});

function onModelChange(modelId) {
  const llm = llmList.value.find(llm=>llm['id'] === modelId);
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
    <Select v-model:value="node.data['llmNodeData']['outputFormat']" :options="nodeConstants.llmOutputFormatOptions"/>
  </FormItem>
  <VariableTable :node-id="node.id"
                 :input-variables="node.data['llmNodeData'].inputVariables"
                 :output-variables="node.data['llmNodeData'].outputVariables"
                 :node-data="node.data"
                 :node="node"/>
</template>

<style scoped>
.prompt-textarea {
  height: 200px;
}
</style>