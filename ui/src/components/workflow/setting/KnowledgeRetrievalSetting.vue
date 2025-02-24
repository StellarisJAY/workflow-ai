<script setup>
import VariableTable from "./VariableTable.vue";
import CommonSetting from "./CommonSetting.vue";
import {Form, FormItem, InputNumber, Slider, Select} from "ant-design-vue";
import {kbSearchTypes} from "../../../api/const.js";

defineProps(['node', 'refOptions', 'kbOptions']);
</script>

<template>
  <CommonSetting :node="node"/>
  <Form>
    <FormItem label="知识库">
      <Select :options="kbOptions" v-model:value="node.data['knowledgeRetrievalNodeData']['kbId']"/>
    </FormItem>
    <FormItem label="搜索类型">
      <Select :options="kbSearchTypes" v-model:value="node.data['knowledgeRetrievalNodeData']['searchType']"/>
    </FormItem>
    <FormItem v-if="node.data['knowledgeRetrievalNodeData']['searchType'] !== 'fulltext'" label="相似度阈值">
      <Slider :min="0.01" :max="0.99" :step="0.01" v-model:value="node.data['knowledgeRetrievalNodeData']['similarityThreshold']"/>
    </FormItem>
    <FormItem label="最大返回文档数量">
      <InputNumber :min="1" v-model:value="node.data['knowledgeRetrievalNodeData']['count']"/>
    </FormItem>
  </Form>
  <VariableTable :node-id="node.id"
                 :has-input="true"
                 :has-output="true"
                 :input-variables="node.data['knowledgeRetrievalNodeData'].inputVariables"
                 :output-variables="node.data['knowledgeRetrievalNodeData'].outputVariables"
                 :ref-options="refOptions"
                 :input-editable="true"
                 :output-editable="false"
                 :allow-add-del-input="false"
                 :allow-add-del-output="false"/>
</template>

<style scoped></style>