<script setup>
import VariableTable from "./VariableTable.vue";

defineProps(['node']);
</script>

<template>
  <VariableTable :input-variables="node.data['startNodeData'].inputVariables"
                 :output-variables="node.data['startNodeData'].outputVariables"
                 :node-id="node.id"
                 :node-data="node.data"
                 :node="node"/>
</template>

<style scoped></style>