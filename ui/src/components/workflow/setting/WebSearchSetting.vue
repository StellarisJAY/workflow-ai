<script setup>
import {FormItem, InputNumber} from "ant-design-vue";
import VariableTable from "./VariableTable.vue";
import CommonSetting from "./CommonSetting.vue";

defineProps(['node']);
</script>

<template>
  <CommonSetting :node="node"/>
  <FormItem label="最大召回数量">
    <InputNumber :min="1" :step="1" v-model:value="node.data['webSearchNodeData']['topN']"/>
  </FormItem>
  <VariableTable :node-id="node.id"
                 :input-variables="node.data['input']"
                 :output-variables="node.data['output']"
                 :node-data="node.data"
                 :node="node"/>
</template>

<style scoped></style>