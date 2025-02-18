<script setup>
import VariableTable from "./VariableTable.vue";

defineProps(['node']);
</script>

<template>
  <VariableTable :input-variables="node.data['startNodeData'].inputVariables"
                 :output-variables="node.data['startNodeData'].outputVariables"
                 :has-input="true"
                 :has-output="false"
                 :allow-ref="false"
                 :node-id="node.id"
                 :input-editable="true"
                 :output-editable="false" :allow-add-del-input="true"/>
</template>

<style scoped></style>