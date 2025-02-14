<script setup>
import {Tag} from "ant-design-vue";

defineProps(['inputVariables', 'outputVariables', 'hasInput', 'hasOutput']);
</script>

<template>
  <span v-if="hasInput">输入: <Tag v-for="variable in inputVariables">{{variable.name}}</Tag></span>
  <span v-if="hasOutput">输出: <Tag v-for="variable in outputVariables">{{variable.name}}</Tag></span>
</template>

<style scoped></style>