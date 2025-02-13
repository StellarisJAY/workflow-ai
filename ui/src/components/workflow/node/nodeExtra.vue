<script setup>
import { DeleteOutlined } from '@ant-design/icons-vue';
import { useVueFlow } from '@vue-flow/core';
import { Button } from 'ant-design-vue';

defineProps(['id', 'type', 'data', 'showControls', 'title']);
const { removeNodes } = useVueFlow();
function removeNode(id) {
  removeNodes(id);
}
</script>

<template>
    <Button type="primary" danger shape="circle" size="small" @click="removeNode(id)">
      <DeleteOutlined/>
    </Button>
</template>

<style scoped></style>