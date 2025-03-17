<script setup>
import {Handle, Position, useVueFlow} from '@vue-flow/core';
import {Card, Tag, List, ListItem} from "ant-design-vue";
import NodeExtra from "./nodeExtra.vue";
import {CheckCircleTwoTone} from "@ant-design/icons-vue";
import {onMounted} from "vue";
import NodeConstants from "../nodeConstants.js";
import NodeUtil from "../../../util/nodeUtil.js";

const props = defineProps(['id', 'type', 'data']);
const {findNode, nodesConnectable, nodesDraggable, getNodes, getEdges} = useVueFlow();
const node = findNode(props.id);

onMounted(()=>{
  addEventListener(NodeConstants.deleteNodeEvent, ev=>onNodeDelete(ev));
  addEventListener(NodeConstants.deleteEdgeEvent, ev=>onEdgeDelete(ev));
});

function opValueToHumanReadable(variable) {
  if (variable.value.type === "literal") {
    return variable.value.content;
  }
  const node = findNode(variable.value.sourceNode);
  if (node) return node.data['name'] + "/" + variable.value.sourceName;
  return "";
}

function onNodeDelete(ev) {
  const nodeId = ev.detail.nodeId;
  const needResetVariableRef = id=>id===nodeId;
  resetVariables(needResetVariableRef);
}

function onEdgeDelete(ev) {
  const prevNodes = NodeUtil.getPrevNodes(props.id, getNodes.value, getEdges.value).map(node=>node.id);
  const needResetVariableRef = id=>!prevNodes.includes(id);
  resetVariables(needResetVariableRef);
}

/**
 * 重置条件列表中的引用
 * @param needResetVariableRef 通过引用节点id判断是否需要重置
 */
function resetVariables(needResetVariableRef) {
  const conditionNodeData = props.data['conditionNodeData'];
  conditionNodeData.branches.forEach(branch=>{
    if (!branch['conditions']) return;
    branch['conditions'].forEach(condition=>{
      if (condition.value1.isRef) {
        const ref = condition.value1.ref.split(".");
        if (ref.length === 2 && needResetVariableRef(ref[0])) {
          condition.value1.ref = "";
          condition.value1.refOption = [];
        }
      }
      if (condition.value2.isRef) {
        const ref = condition.value2.ref.split(".");
        if (ref.length === 2 && needResetVariableRef(ref[0])) {
          condition.value2.ref = "";
          condition.value2.refOption = [];
        }
      }
    });
  });
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
          <Tag>{{opValueToHumanReadable(condition.value1)}}</Tag>
          <Tag>{{condition.op}}</Tag>
          <Tag>{{opValueToHumanReadable(condition.value2)}}</Tag>
        </ListItem>
        <ListItem v-if="i < branch.conditions.length-1" style="margin:2px; padding: 2px">
          {{branch['connector']}}
        </ListItem>
      </List>
    </Card>
  </Card>
</template>

<style scoped></style>