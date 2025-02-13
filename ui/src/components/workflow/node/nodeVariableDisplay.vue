<script setup>
import {Tag} from "ant-design-vue";

defineProps(['inputVariables', 'outputVariables', 'hasInput', 'hasOutput']);
</script>

<template>
  <h4 v-if="hasInput">输入: </h4>
  <Tag v-if="hasInput" v-for="variable in inputVariables">{{variable.name}}</Tag>
  <h4 v-if="hasOutput">输出: </h4>
  <Tag v-if="hasOutput" v-for="variable in outputVariables">{{variable.name}}</Tag>
</template>

<style scoped></style>