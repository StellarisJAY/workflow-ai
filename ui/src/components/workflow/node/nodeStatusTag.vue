<script setup>
import { CheckCircleOutlined, ClockCircleOutlined, SyncOutlined, CloseCircleOutlined }
  from '@ant-design/icons-vue';
import { Tag } from 'ant-design-vue';

const props = defineProps(['status']);
function getStatusTagColor() {
  switch (props.status.id) {
    case 0: return 'default';
    case 1: return 'processing';
    case 2: return 'success';
    case 3: return 'error';
  }
  return 'default';
}
</script>

<template>
  <Tag :color="getStatusTagColor()">
    <template #icon>
      <ClockCircleOutlined v-if="status.id === 0" />
      <SyncOutlined v-else-if="status.id === 1" :spin="true"/>
      <CheckCircleOutlined v-else-if="status.id === 2"/>
      <CloseCircleOutlined v-else-if="status.id === 3" />
    </template>
    {{status.text}}
  </Tag>
</template>

<style scoped></style>