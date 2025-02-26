<script setup>
import VariableTable from "./VariableTable.vue";
import CommonSetting from "./CommonSetting.vue";

defineProps(['node']);
</script>

<template>
  <CommonSetting :node="node"/>
  <VariableTable :node-id="node.id"
                 :node-data="node.data"
                 :node="node"
                 :input-variables="node.data['crawlerNodeData'].inputVariables"
                 :output-variables="node.data['crawlerNodeData'].outputVariables"/>
</template>

<style scoped></style>