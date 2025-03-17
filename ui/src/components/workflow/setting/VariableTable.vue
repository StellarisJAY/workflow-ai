<script setup>
import {Input, List, ListItem, Button, Select, Cascader} from "ant-design-vue";
const props = defineProps(['inputVariables', 'outputVariables','nodeId', "nodeData", "node"]);
import {DeleteFilled} from "@ant-design/icons-vue";
import {onMounted, ref} from "vue";
import NodeUtil from "../../../util/nodeUtil.js";

const refOptions = ref([]);

onMounted(()=>{
  refOptions.value = NodeUtil.getPrevNodesOutputs(props.nodeId);
  if (!props.inputVariables) return;
  props.inputVariables.forEach(variable => {
    if (variable.value.type === "ref") {
      variable['refOption'] = [variable.value.sourceNode, variable.value.sourceName];
    }
  });
});

function onRefOptionChange(variable, ev) {
  variable.value['sourceNode'] = ev[0];
  variable.value['sourceName'] = ev[1];
  // 引用发生变化，将当前变量类型改为被引用变量类型
  const srcNode = refOptions.value.find(option=>option.value === ev[0]);
  if (srcNode) {
    const srcVar = srcNode.children.find(child=>child.label === ev[1]);
    if (srcVar) {
      variable.type = srcVar.type;
    }
  }
}

function addVariable(target) {
  target.push({name: "variable", type: "string", value: {
      type: "ref",
      content: "",
      sourceNode: "",
      sourceName: "",
    }});
}

function removeVariable(target, name) {
  const idx = target.findIndex(item=>item.name === name);
  if (idx > -1) {
    target.splice(idx, 1);
  }
}

function getTypeOptions(variable) {
  let types;
  if (variable['allowedTypes']) {
    types = variable['allowedTypes'];
  }else {
    types = props.nodeData['defaultAllowVarTypes'];
  }
  const typeOptions = [];
  types.forEach(item => {
    typeOptions.push({label: item, value: item});
  });
  return typeOptions;
}
</script>

<template>
  <h4 v-if="inputVariables">输入变量</h4>
  <List>
    <ListItem v-for="variable in inputVariables">
      <Input v-model:value="variable.name"
             size="small"
             placeholder="变量名"
             :disabled="variable['required']"/>
      <!--引用-->
      <Cascader v-if="node.type !== 'start'" size="small"
                v-model:value="variable['refOption']"
                :options="refOptions"
                @change="ev=>onRefOptionChange(variable, ev)"/>
      <Select v-else v-model:value="variable.type"
              :options="getTypeOptions(variable)"
              size="small"
              :disabled="variable['fixed']"/>
      <Button size="small"
              v-if="!variable['required'] && !variable['fixed']"
              @click="removeVariable(inputVariables, variable.name)"><DeleteFilled/></Button>
    </ListItem>
    <ListItem>
      <Button v-if="nodeData['allowAddInputVar']"
              @click="addVariable(inputVariables)"
              size="mini">添加</Button>
    </ListItem>
  </List>
  <h4 v-if="outputVariables">输出变量</h4>
  <List>
    <ListItem v-for="variable in outputVariables">
      <Input v-model:value="variable.name"
             size="small" placeholder="变量名"
             :disabled="true"/>
      <Select v-model:value="variable.type"
              :options="getTypeOptions(variable)"
              size="small"
              :disabled="true"/>
    </ListItem>
  </List>
</template>

<style scoped></style>