<script setup>
import CommonSetting from "./CommonSetting.vue";
import {Card, List, ListItem, Input, Select, Cascader, Button} from "ant-design-vue";
import {DeleteOutlined} from "@ant-design/icons-vue";
import types from "../types.js";
import {useVueFlow} from "@vue-flow/core";

const props = defineProps(['node', 'refOptions']);
const {getEdges, removeEdges} = useVueFlow();

const varTypeOptions = [
  {label: "引用", value: "ref"},
  {label: "string", value: "string"},
];

const operators = [
  {label: "等于", value: "=="},
  {label: "不等于", value: "!="},
  {label: "大于", value: ">"},
  {label: "小于", value: "<"},
  {label: "大于等于", value: ">="},
  {label: "小于等于", value: "<="},
];

function addCondition(branch) {
  branch.conditions.push({
    value1: {type: "string", value: ""},
    op: "==",
    value2: {type: "string", value: ""},
  });
}

function addBranch() {
  const branches = props.node.data['conditionNodeData']['branches'];
  branches.splice(branches.length-1, 0, types.createBranch());
}

function removeBranch(idx) {
  if (idx === 0) return;
  const conditionNode = props.node.data['conditionNodeData'];
  const branches = conditionNode['branches'];
  const deletedBranch = branches.splice(idx, 1)[0];
  const edge = getEdges.value.find((edge) => edge.targetHandle === deletedBranch.handle);
  if (edge) {
    removeEdges(edge.id);
  }
}

function removeCondition(branch, idx, branchIdx) {
  branch.conditions.splice(idx, 1);
  if (branch.conditions.length === 0) {
    removeBranch(branchIdx);
  }
}

function onValueTypeChange(ev) {

}

</script>

<template>
  <CommonSetting :node="node"/>
  <div v-for="(branch, branchIdx) in node.data['conditionNodeData']['branches']" :key="branch.handle">
    <Card v-if="branchIdx < node.data['conditionNodeData']['branches'].length-1">
      <template #title>
        <h5 v-if="branchIdx === 0">如果</h5>
        <h5 v-else>否则如果</h5>
      </template>
      <template #extra>
        <Button @click="removeBranch(branchIdx)">删除</Button>
      </template>
      <List v-if="branchIdx < node.data['conditionNodeData']['branches'].length-1">
        <ListItem v-for="(condition, idx) in branch['conditions']">
          <!--操作数1-->
          <Select :options="varTypeOptions" v-model:value="condition.value1.type"
                  @change="_=>{condition.value1.value = '';}"/>
          <Input v-if="condition.value1.type==='string'" v-model:value="condition.value1.value"/>
          <Cascader v-else-if="condition.value1.type==='ref'"
                    v-model:value="condition.value1.refOption"
                    @change="ev=>{condition.value1.value = ev.join('.');}"
                    :options="refOptions"/>
          <!--符号-->
          <Select :options="operators" v-model:value="condition.op"></Select>
          <!--操作数2-->
          <Select :options="varTypeOptions" v-model:value="condition.value2.type"
                  @change="_=>{condition.value1.value = '';}"/>
          <Input v-if="condition.value2.type==='string'" v-model:value="condition.value2.value"/>
          <Cascader v-else-if="condition.value2.type==='ref'"
                    v-model:value="condition.value2.refOption"
                    @change="ev=>{condition.value2.value = ev.join('.');}"
                    :options="refOptions"/>

          <Button @click="removeCondition(branch, idx, branchIdx)"><DeleteOutlined/></Button>
        </ListItem>
        <ListItem><Button @click="addCondition(branch)">添加</Button></ListItem>
      </List>
    </Card>
  </div>
  <Button @click="addBranch">添加分支</Button>

</template>

<style scoped></style>