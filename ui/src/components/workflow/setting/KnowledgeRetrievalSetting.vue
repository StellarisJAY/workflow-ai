<script setup>
import VariableTable from "./VariableTable.vue";
import CommonSetting from "./CommonSetting.vue";
import {Form, FormItem, InputNumber, Slider, Select} from "ant-design-vue";
import {kbSearchTypes} from "../../../api/const.js";
import {onMounted, ref} from "vue";
import knowledgeBaseAPI from "../../../api/knowledgeBase.js";

defineProps(['node']);

const kbList = ref([]);
const kbOptions = ref([]);

onMounted(()=>{
  knowledgeBaseAPI.list({paged: false}).then(resp=>{
    kbList.value = resp.data;
    const options = [];
    kbList.value.forEach(item=>{
      options.push({label: item.name, value: item.id});
    });
    kbOptions.value = options;
  });
});

</script>

<template>
  <CommonSetting :node="node"/>
  <Form>
    <FormItem label="知识库">
      <Select :options="kbOptions" v-model:value="node.data['retrieveKnowledgeBaseNodeData']['kbId']"/>
    </FormItem>
    <FormItem label="搜索类型">
      <Select :options="kbSearchTypes" v-model:value="node.data['retrieveKnowledgeBaseNodeData']['searchType']"/>
    </FormItem>
    <FormItem v-if="node.data['retrieveKnowledgeBaseNodeData']['searchType'] !== 'fulltext'" label="相似度阈值">
      <Slider :min="0.01" :max="0.99" :step="0.01" v-model:value="node.data['retrieveKnowledgeBaseNodeData']['similarityThreshold']"/>
    </FormItem>
    <FormItem label="最大返回文档数量">
      <InputNumber :min="1" v-model:value="node.data['retrieveKnowledgeBaseNodeData']['count']"/>
    </FormItem>
  </Form>
  <VariableTable :node-id="node.id"
                 :node-data="node.data"
                 :node="node"
                 :input-variables="node.data['retrieveKnowledgeBaseNodeData'].inputVariables"
                 :output-variables="node.data['retrieveKnowledgeBaseNodeData'].outputVariables"/>
</template>

<style scoped></style>