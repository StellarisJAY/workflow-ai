<script setup>
import {List, ListItem, Timeline, TimelineItem} from "ant-design-vue";

defineProps(['outputs'])
</script>

<template>
  <Timeline mode="right" style="left: 0">
    <TimelineItem v-for="output in outputs" :label="output['nodeName']">
      {{output['output']}}
    </TimelineItem>
  </Timeline>
</template>

<style scoped></style>