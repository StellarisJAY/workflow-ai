<script setup>
import {Button, Cascader, Input, List, ListItem, Select} from "ant-design-vue";
import {DeleteFilled} from "@ant-design/icons-vue";

defineProps(['node', 'outputVariables', 'refOptions']);

const typeOptions = [
  {label: "string", value: "string"},
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
  <h4>输出变量</h4>
  <List>
    <ListItem v-for="variable in outputVariables">
      <Input v-model:value="variable.name" size="small" placeholder="变量名"></Input>
      <Select v-model:value="variable.type" :options="typeOptions" size="small"></Select>
      <Input v-if="variable.type === 'string'" v-model:value="variable.value" size="small" placeholder="值"></Input>
      <Cascader v-else-if="variable.type === 'ref'"
                :options="refOptions" size="small"
                @change="ev=>onRefOptionChange(variable, ev)"></Cascader>
      <Button size="small"
              @click="removeVariable(outputVariables, variable.name)"><DeleteFilled/></Button>
    </ListItem>
    <ListItem>
      <Button @click="addVariable(outputVariables)" size="mini">添加</Button>
    </ListItem>
  </List>
</template>

<style scoped></style>