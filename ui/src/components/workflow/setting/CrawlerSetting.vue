<script setup>
import VariableTable from "./VariableTable.vue";
import CommonSetting from "./CommonSetting.vue";

defineProps(['node']);
</script>

<template>
  <CommonSetting :node="node"/>
  <VariableTable :node-id="node.id"
                 :input-variables="node.data['input']"
                 :output-variables="node.data['output']"
                 :node-data="node.data"
                 :node="node"/>
</template>

<style scoped></style>