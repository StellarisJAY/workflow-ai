<script setup>
import VariableTable from "./VariableTable.vue";

defineProps(['node', 'refOptions']);
</script>

<template>
  <VariableTable :node-id="node.id"
                 :has-input="true"
                 :has-output="true"
                 :input-variables="node.data['crawlerNodeData'].inputVariables"
                 :output-variables="node.data['crawlerNodeData'].outputVariables"
                 :ref-options="refOptions"
                 :input-editable="true"
                 :output-editable="false"/>
</template>

<style scoped></style>