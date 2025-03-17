<script setup>
import VariableTable from "./VariableTable.vue";

const props = defineProps(['node', 'outputVariables']);
</script>

<template>
  <VariableTable :node="node" :node-data="node.data" :input-variables="node.data.input" :node-id="node.id" />
</template>

<style scoped></style>