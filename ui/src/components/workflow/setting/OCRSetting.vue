<script setup>
import {FormItem, Select} from "ant-design-vue";
import VariableTable from "./VariableTable.vue";
import CommonSetting from "./CommonSetting.vue";
import {onMounted, ref} from "vue";
import llmAPI from "../../../api/llm.js";
const props = defineProps(['node']);

const llmList = ref([]);
const llmOptions = ref([]);

onMounted(()=>{
  llmAPI.listModels({paged: false, modelType: "image_understanding"}).then(resp=>{
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
    props.node.data['ocrNodeData']['modelName'] = llm.name;
  }
}
</script>

<template>
  <CommonSetting :node="node"/>
  <FormItem label="模型">
    <Select v-model:value="node.data['ocrNodeData']['modelId']" :options="llmOptions" @change="onModelChange"></Select>
  </FormItem>
  <VariableTable :node-id="node.id"
                 :input-variables="node.data['input']"
                 :output-variables="node.data['output']"
                 :node-data="node.data"
                 :node="node"/>
</template>

<style scoped>
</style>