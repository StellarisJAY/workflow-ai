<script setup>
import {Input, FormItem} from "ant-design-vue";

defineProps(['node']);
</script>

<template>
  <FormItem label="节点名称">
    <Input v-model:value="node.data['name']" placeholder="节点名称"/>
  </FormItem>
</template>

<style scoped></style>