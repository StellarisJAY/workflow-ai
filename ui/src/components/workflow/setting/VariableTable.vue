<script setup>
import {Input, List, ListItem, Button, Select, Cascader} from "ant-design-vue";
const props = defineProps(['inputVariables', 'outputVariables','nodeId', "nodeData", "node"]);
import {DeleteFilled} from "@ant-design/icons-vue";
import {onMounted, ref} from "vue";
import NodeUtil from "../../../util/nodeUtil.js";

const refOptions = ref([]);
onMounted(()=>{
  refOptions.value = NodeUtil.getPrevNodesOutputs(props.nodeId);
  props.inputVariables.forEach(variable => {
    if (variable.ref) {
      variable['refOption'] = variable['ref'].split('.');
    }
  });
});

function onRefOptionChange(variable, ev) {
  variable['ref'] = ev[0] + "." + ev[1];
}

function addVariable(target) {
  let allowRef = true;
  const nodeType = props.node.type;
  if (nodeType === 'start') allowRef = false;
  target.push({name: "variable", value: "", type: "string", allowRef: allowRef});
}

function addOutputVariable(target) {
  target.push({name: "variable", value: "", type: "string", allowRef: false});
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

function getValueOptions(variable) {
  const options = [{label: "值", value: "value"}];
  if (variable['allowRef']) {
    options.push({label: "引用", value: "ref"});
  }
  return options;
}

function getValueOption(variable) {
  return variable['isRef'] ? "ref":"value";
}

function onValueOptionChange(variable, ev) {
  variable['isRef'] = ev === 'ref';
}
</script>

<template>
  <h4>输入变量</h4>
  <List>
    <ListItem v-for="variable in inputVariables">
      <Input v-model:value="variable.name"
             size="small"
             placeholder="变量名"
             :disabled="variable['required']"/>
      <!--值选择-->
      <Select :value="getValueOption(variable)" size="small"
              :options="getValueOptions(variable)"
              @change="ev=>onValueOptionChange(variable, ev)"/>
      <div v-if="!variable.isRef">
        <!--类型选择-->
        <Select v-model:value="variable.type"
                :options="getTypeOptions(variable)"
                size="small"
                :disabled="variable['fixed']"/>
        <!--字符串-->
        <Input v-model:value="variable['value']" size="small"/>
      </div>

      <!--引用-->
      <Cascader v-else size="small"
                v-model:value="variable['refOption']"
                :options="refOptions"
                @change="ev=>onRefOptionChange(variable, ev)"/>

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
  <h4>输出变量</h4>
  <List>
    <ListItem v-for="variable in outputVariables">
      <Input v-model:value="variable.name"
             size="small" placeholder="变量名"
             :disabled="variable['required']"/>
      <Select v-model:value="variable.type"
              :options="getTypeOptions(variable)"
              size="small"
              :disabled="variable['fixed']"/>
      <Button size="small"
              @click="removeVariable(outputVariables, variable.name)"
              v-if="!variable['required']&&!variable['fixed']"><DeleteFilled/></Button>
    </ListItem>
    <ListItem>
      <Button @click="addOutputVariable(outputVariables)"
              size="mini"
              v-if="nodeData['allowAddOutputVar']">添加</Button>
    </ListItem>
  </List>
</template>

<style scoped></style>