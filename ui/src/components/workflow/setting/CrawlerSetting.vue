<script setup>
import VariableTable from "./VariableTable.vue";
import CommonSetting from "./CommonSetting.vue";

defineProps(['node', 'refOptions']);
</script>

<template>
  <CommonSetting :node="node"/>
  <VariableTable :node-id="node.id"
                 :has-input="true"
                 :has-output="true"
                 :input-variables="node.data['crawlerNodeData'].inputVariables"
                 :output-variables="node.data['crawlerNodeData'].outputVariables"
                 :ref-options="refOptions"
                 :input-editable="true"
                 :output-editable="false"
                 :allow-add-del-input="false"
                 :allow-add-del-output="false"/>
</template>

<style scoped></style>