<script setup>
import {Row, Col, Input, List, ListItem, Button, FormItem, Select, Cascader} from "ant-design-vue";
const props = defineProps(['inputVariables', 'outputVariables', 'hasInput', 'hasOutput',
  'allowRef', 'nodeId', 'refOptions']);
import {DeleteFilled} from "@ant-design/icons-vue";

const typeOptions = [
  {label: "string", value: "string"},
  {label: "文件", value: "file"},
  {label: "引用", value: "ref"},
];
function onRefOptionChange(variable, ev) {
  variable.value = ev[0] + "." + ev[1];
}

function addVariable(target) {
  target.push({name: "", value: "", type: "string"});
}

function removeVariable(target, name) {
  const idx = target.findIndex(item=>item.name !== name);
  target.splice(idx, 1);
}
</script>

<template>
  <Row v-if="hasInput">
    <Col :span="8">
      <h4>输入变量</h4>
    </Col>
    <Col :span="4">
      <Button @click="addVariable(inputVariables)">添加</Button>
    </Col>
  </Row>
  <List v-if="hasInput">
    <ListItem v-for="variable in inputVariables">
      <Row>
       <Col :span="8">
         <Input v-model:value="variable.name" size="small"></Input>
       </Col>
        <Col :span="7">
          <Select v-model:value="variable.type" :options="typeOptions" size="small"></Select>
        </Col>
        <Col :span="8">
          <Input v-if="variable.type === 'string'" v-model:value="variable.value" size="small"></Input>
          <Cascader v-else-if="variable.type === 'ref'"
                    :options="refOptions" size="small"
                    @change="ev=>onRefOptionChange(variable, ev)"></Cascader>
        </Col>
        <Col :span="1">
          <Button size="mini" style="width: 20px; height: 20px" @click="removeVariable(inputVariables, variable.name)"><DeleteFilled/></Button>
        </Col>
      </Row>
    </ListItem>
  </List>
  <h4 v-if="hasOutput">输出变量</h4>
  <List v-if="hasOutput">
    <ListItem v-for="variable in outputVariables">
      <Row>
        <Col :span="8">
          <Input v-model:value="variable.name" size="small"></Input>
        </Col>
        <Col :span="7">
          <Select v-model:value="variable.type" :options="typeOptions" size="small"></Select>
        </Col>
        <Col :span="8">
          <Input v-if="variable.type === 'string'" v-model:value="variable.value" size="small"></Input>
          <Cascader v-else-if="variable.type === 'ref'"
                    :options="refOptions" size="small"
                    @change="ev=>onRefOptionChange(variable, ev)"></Cascader>
        </Col>
        <Col :span="1">
          <Button size="mini" style="width: 20px; height: 20px" @click="removeVariable(inputVariables, variable.name)"><DeleteFilled/></Button>
        </Col>
      </Row>
    </ListItem>
  </List>
</template>

<style scoped></style>