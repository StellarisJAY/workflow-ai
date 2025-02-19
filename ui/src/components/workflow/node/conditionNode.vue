<script setup>
import {Handle, Position, useVueFlow} from '@vue-flow/core';
import {Card, Tag, List, ListItem} from "ant-design-vue";
import NodeExtra from "./nodeExtra.vue";
import {CheckCircleTwoTone} from "@ant-design/icons-vue";

const props = defineProps(['id', 'type', 'data']);
const {findNode, nodesConnectable, nodesDraggable} = useVueFlow();
const node = findNode(props.id);

function opValueToHumanReadable(type, val) {
  if (type === 'ref') {
    const parts = val.split('.');
    if (parts.length > 1) {
      const node = findNode(parts[0]);
      return node.data['name'] + "/" + parts[1];
    }
  }
  return val;
}
</script>

<template>
  <Card :title="data['name']" :hoverable="true" :body-style="{padding:'10px'}" :head-style="{padding:'10px'}">
    <Handle type="target" :position="Position.Left"></Handle>
    <template #extra>
      <node-extra :id="id" :type="type" :data="data" :status="node.status" :editable="true"/>
    </template>
    <Card v-for="(branch, branchIdx) in data['conditionNodeData']['branches']"
          :key="branch.handle"
          :body-style="{padding:'2px'}" :head-style="{padding:'5px', margin:'0', fontSize:'10px'}" >
      <template #title>
        <p v-if="branchIdx === 0">如果</p>
        <p v-else-if="branchIdx < data['conditionNodeData']['branches'].length-1">否则如果</p>
        <p v-else>否则</p>
      </template>

      <template #extra>
        <CheckCircleTwoTone two-tone-color="#00ff00" v-if="!nodesConnectable && !nodesDraggable && branch['success']" />
      </template>
      <Handle :id="branch['handle']" type="source" :position="Position.Right"></Handle>
      <List v-for="(condition, i) in branch.conditions">
        <ListItem style="margin:2px; padding: 2px">
          <Tag>{{opValueToHumanReadable(condition.value1.type, condition.value1.value)}}</Tag>
          <Tag>{{condition.op}}</Tag>
          <Tag>{{opValueToHumanReadable(condition.value2.type, condition.value2.value)}}</Tag>
        </ListItem>
        <ListItem v-if="i < branch.conditions.length-1" style="margin:2px; padding: 2px">
          {{branch['connector']}}
        </ListItem>
      </List>
    </Card>
  </Card>
</template>

<style scoped></style>