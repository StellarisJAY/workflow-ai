<script setup>
import {Button, Cascader, Input, List, ListItem, Select} from "ant-design-vue";
import {DeleteFilled} from "@ant-design/icons-vue";
import {onMounted, ref} from "vue";
import NodeUtil from "../../../util/nodeUtil.js";

const props = defineProps(['node', 'outputVariables']);
const refOptions = ref([]);

onMounted(()=>{
  refOptions.value = NodeUtil.getPrevNodesOutputs(props.node.id);
  props.outputVariables.forEach(variable => {
    if (variable.isRef) {
      variable['refOption'] = variable.ref.split('.');
    }
  });
});

function onRefOptionChange(variable, ev) {
  variable.ref = ev[0] + "." + ev[1];
  variable.isRef = true;
  const node = refOptions.value.find(option=>option["value"] === ev[0]);
  if (node) {
    const child = node["children"].find(child=>child["label"] === ev[1]);
    if (child) {
      variable.type = child["type"];
    }
  }
}

function addVariable(target) {
  target.push({name: "variable", value: "", isRef: true, ref: "", type: "string"});
}

function removeVariable(target, name) {
  const idx = target.findIndex(item=>item.name === name);
  target.splice(idx, 1);
}
</script>

<template>
  <h4>输出变量</h4>
  <List>
    <ListItem v-for="variable in outputVariables">
      <Input v-model:value="variable.name" size="small" placeholder="变量名"></Input>
      <Cascader :options="refOptions" size="small"
                @change="ev=>onRefOptionChange(variable, ev)" v-model:value="variable.refOption"></Cascader>
      <Button size="small"
              @click="removeVariable(outputVariables, variable.name)"><DeleteFilled/></Button>
    </ListItem>
    <ListItem>
      <Button @click="addVariable(outputVariables)" size="mini">添加</Button>
    </ListItem>
  </List>
</template>

<style scoped></style>