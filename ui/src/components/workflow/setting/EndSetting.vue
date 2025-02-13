<script setup>
import VariableTable from "./VariableTable.vue";

defineProps(['node']);
</script>

<template>
  <VariableTable :output-variables="node.data.outputVariables"
                 :has-input="false"
                 :has-output="true"
                 :allow-ref="true"
                 :node-id="node.id"/>
</template>

<style scoped></style>