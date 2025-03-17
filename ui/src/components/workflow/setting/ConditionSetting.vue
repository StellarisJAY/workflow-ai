<script setup>
import CommonSetting from "./CommonSetting.vue";
import {Card, List, ListItem, Input, Select, Cascader, Button, InputNumber} from "ant-design-vue";
import {DeleteOutlined} from "@ant-design/icons-vue";
import nodeConstants from "../nodeConstants.js";
import {useVueFlow} from "@vue-flow/core";
import {onMounted, ref} from "vue";
import NodeUtil from "../../../util/nodeUtil.js";
import NodeConstants from "../nodeConstants.js";

const props = defineProps(['node']);
const {getEdges, removeEdges} = useVueFlow();

const varTypeOptions = [
  {label: "string", value: "string"},
  {label: "number", value: "number"},
];
const varValueOptions = [
  {label: "值", value: "value"},
  {label: "引用", value: "ref"},
];

const numberOperators = [
  {label: "等于", value: "=="},
  {label: "不等于", value: "!="},
  {label: "大于", value: ">"},
  {label: "小于", value: "<"},
  {label: "大于等于", value: ">="},
  {label: "小于等于", value: "<="},
];

const stringOperators = [
  {label: "等于", value: "=="},
  {label: "不等于", value: "!="},
  {label: "包含", value: "contains"},
  {label: "不包含", value: "!contains"},
  {label: "为空", value: "empty"},
  {label: "不为空", value: "!empty"},
];

const arrayOperators = [
  {label: "为空", value: "empty"},
  {label: "不为空", value: "!empty"},
];

const refOptions = ref([]);

onMounted(()=>{
  refOptions.value = NodeUtil.getPrevNodesOutputs(props.node.id);
  const nodeData = props.node.data["conditionNodeData"];
  nodeData["branches"].forEach(branch=>{
    if (branch["handle"] === "else") return;
    branch['conditions'].forEach(condition=>{
      const value1 = condition.value1;
      const value2 = condition.value2;
      if (value1.value.type === "ref") {
        value1["refOption"] = [value1.value.sourceNode, value1.value.sourceName];
      }
      if (value2.value.type === "ref") {
        value2["refOption"] = [value2.value.sourceNode, value2.value.sourceName];
      }
    });
  });
});

function getOperatorsOfType(type) {
  if (type === "number") {
    return numberOperators;
  }
  if (type === "string") {
    return stringOperators;
  }
  if (type === "array_str" || type === "array_num") {
    return arrayOperators;
  }
  return [];
}

function addCondition(branch) {
  branch.conditions.push({
    value1: {type: "string", value: {type: "ref", content: "", sourceNode: "", sourceName: ""}},
    op: "==",
    value2: {type: "string", value: {type: "literal", content: "", sourceNode: "", sourceName: ""}},
  });
}

function isValue2Needed(op) {
  return op !== "empty" && op !== "!empty"
}

function addBranch() {
  const branches = props.node.data['conditionNodeData']['branches'];
  branches.splice(branches.length-1, 0, nodeConstants.createBranch());
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

function onValueRefChange(value1, value2, ev) {
  value1.value.sourceNode = ev[0];
  value1.value.sourceName = ev[1];
  // 引用发生变化，将当前变量类型改为被引用变量类型
  const srcNode = refOptions.value.find(option=>option.value === ev[0]);
  if (srcNode) {
    const srcVar = srcNode.children.find(child=>child.label === ev[1]);
    if (srcVar) {
      value1.type = srcVar.type;
      value2.type = srcVar.type;
    }
  }
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
          <Cascader
                    v-model:value="condition.value1.refOption"
                    @change="ev=>onValueRefChange(condition.value1, condition.value2, ev)"
                    :options="refOptions"/>

          <!--符号-->
          <Select :options="getOperatorsOfType(condition.value1.type)" v-model:value="condition.op"></Select>
          <div v-if="isValue2Needed(condition.op)">
            <Input v-if="condition.value1.type === 'string'" v-model:value="condition.value2.value.content"/>
            <InputNumber v-else-if="condition.value1.type === 'number'"
                         v-model:value="condition.value2.value.content" :string-mode="true"/>
          </div>
          <Button @click="removeCondition(branch, idx, branchIdx)"><DeleteOutlined/></Button>
        </ListItem>
        <ListItem><Button @click="addCondition(branch)">添加</Button></ListItem>
      </List>
    </Card>
  </div>
  <Button @click="addBranch">添加分支</Button>

</template>

<style scoped></style>