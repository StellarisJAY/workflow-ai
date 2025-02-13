<script setup>
import { markRaw, ref } from 'vue'
import { Panel, ConnectionMode, Position, useVueFlow, VueFlow, MarkerType } from '@vue-flow/core'
import { Background } from '@vue-flow/background';
import LLMNode from './node/llmNode.vue';
import CustomEdge from './customEdge.vue';
import StartNode from './node/startNode.vue';
import EndNode from './node/endNode.vue';
import knowledgeRetrievalNode from './node/knowledgeRetrivalNode.vue';
import knowledgeWriteNode from './node/knowledgeWriteNode.vue';
import conditionNode from './node/conditionNode.vue';
import { Button, PageHeader, Select, Modal, Drawer } from 'ant-design-vue';
import llmAPI from '../../api/llm';
import templateAPI from '../../api/template.js';
import LlmSetting from './setting/LLMSetting.vue';
import StartSetting from "./setting/StartSetting.vue";
import NodeUtil from "../../util/nodeUtil.js";

const nodeTypeOptions = ref([
	{ value: "llm", label: "大模型", description: "使用提示词和变量让大模型生成内容" },
	{ value: "knowledgeRetrieval", label: "知识库检索", description: "从知识库检索数据" },
	{ value: "knowledgeWrite", label: "知识库写入", description: "将数据写入知识库" },
	{ value: "condition", label: "条件", description: "条件分支" }
]);

const selectNodeType = ref("llm");

const nodeTypes = ref({
	llm: markRaw(LLMNode),
	start: markRaw(StartNode),
	end: markRaw(EndNode),
	condition: markRaw(conditionNode),
	knowledgeRetrieval: markRaw(knowledgeRetrievalNode),
	knowledgeWrite: markRaw(knowledgeWriteNode)
});

const edgeTypes = ref({
	custom: markRaw(CustomEdge),
});

const nodes = ref([
	{
		id: "start",
		type: "start",
		position: { x: 100, y: 100 },
		data: {
      label: "开始",
      inputVariables: [
        {name: "input", type: "string", "value": ""}
      ]
    }
	},
	{
		id: "end",
		type: "end",
		position: { x: 200, y: 100 },
		data: {
      label: "结束",
      outputVariables: [
        {name: "output", type: "string", value: ""}
      ]
    }
	}
]);
const edges = ref([]);
const newNodeModalOpen = ref(false);
const { onNodeClick, onConnect, onEdgesChange, onNodeDragStop, 
	onNodesChange, onNodeMouseEnter, onNodeMouseLeave} = useVueFlow();

const llmDrawerOpen = ref(false);
const knowledgeRetrievalDrawerOpen = ref(false);
const knowledgeWriteDrawerOpen = ref(false);
const startDrawerOpen = ref(false);
const endDrawerOpen = ref(false);
const currentSettingNode = ref({});

const llmList = ref([]);
const llmOptions = ref([]);
const settingRefOptions = ref([]);

// 点击节点，弹出侧边设置
onNodeClick(event => {
	currentSettingNode.value = event.node;
  settingRefOptions.value = getPrevNodesOutputs();
	switch (event.node.type) {
		case "llm": prepareLLMOptions(); break;
		case "knowledgeRetrieval": knowledgeRetrievalDrawerOpen.value = true; break;
		case "knowledgeWrite": knowledgeWriteDrawerOpen.value = true; break;
    case "start": startDrawerOpen.value = true; break;
    case "end": endDrawerOpen.value = true; break;
	}
});
// 节点连线事件，添加edge
onConnect(event => {
  const exist = edges.value.find(e=>{return e['source'] === event['source'] && e['target'] === event['target']});
  if (exist) {
    return;
  }
	edges.value.push({
		id: crypto.randomUUID(),
		type: "custom",
		source: event.source,
		target: event.target,
		markerStart: MarkerType.ArrowClosed,
	});
});
// 连线断开
onEdgesChange(ev => {
	ev.forEach(e => {
		if (e.type === "remove") {
			removeEdge(e.id);
		}
	});
});
// 节点删除事件
onNodesChange(ev=>{
	ev.forEach(e=>{
		if (e.type === "remove") {
			removeNode(e.id);
		}
	})
});
// 节点移动事件
onNodeDragStop(ev => {
	nodes.value.find(n => n.id === ev.node.id).position = ev.node.position;
});

// 添加节点
function addNode(nodeType) {
	const id = crypto.randomUUID();
	const node = {
		type: nodeType,
		id: id,
		position: { x: 300, y: 200 },
		data: {}
	};
	switch (nodeType) {
		case "llm": node.data = initLLMNodeData(); break;
		case "knowledgeRetrival": break;
		case "knowledgeWrite": break;
		case "condition": break;
	}
	nodes.value.push(node);
}
// 删除节点
function removeNode(id) {
	nodes.value = nodes.value.filter(n => n.id !== id);
}
// 删除连线
function removeEdge(id) {
	edges.value = edges.value.filter(e=>e.id !== id);
}
// 获取流程模板JSON
function getJSON() {
  nodes.value.forEach((node)=>{
    const temp = node.data;
    node.data = null;
    switch (node.type) {
      case "start": node.data = {startNodeData: temp}; break;
      case "end": node.data = {endNodeData: temp}; break;
      case "llm": node.data = {llmNodeData: temp}; break;
    }
  });
	return JSON.stringify({
		nodes: nodes.value,
		edges: edges.value,
	});
}

function closeNewNodeModal() {
	newNodeModalOpen.value = false;
}

function openNewNodeModal() {
	newNodeModalOpen.value = true;
}

function newNodeConfirm() {
	addNode(selectNodeType.value);
	closeNewNodeModal();
}
// 上传模板
function saveTemplate() {
	const data = getJSON();
  templateAPI.createTemplate({
    name: "测试模板1",
    data: data,
  });
}
// 初始化大模型节点数据
function initLLMNodeData() {
	return {
    label: "大模型",
		inputVariables: [{name:"input",type:"string",value:""}],
		outputVariables: [{name: "output", type: "string", value: ""}],
	};
}
// 获取前驱节点的输出变量列表
function getPrevNodesOutputs() {
  const prevNodes = NodeUtil.getPrevNodes(currentSettingNode.value.id, nodes.value, edges.value);
  let options = [];
  prevNodes.forEach(node=>{
    if (!node.data) return;
    let outputVariables = node.data.outputVariables;
    if (node.type === "start") outputVariables = node.data.inputVariables;
    if (outputVariables) {
      let option = {
        label: node.data.label,
        value: node.id,
        children: []
      };
      outputVariables.forEach(variable=>{
        option.children.push({label: variable.name, value: variable.name});
      });
      options.push(option);
    }
  });
  return options;
}

function prepareLLMOptions() {
  llmAPI.listModels({}).then(resp=>{
    console.log(resp.data);
    llmList.value = resp.data;
    const options = [];
    llmList.value.forEach(item=>{
      options.push({label: item.name, value: item.id});
    });
    llmOptions.value = options;
    llmDrawerOpen.value = true;
  });
}
</script>

<template>
	<page-header title="配置流程" style="border: 1px solid rgb(235, 237, 240); height: 10vh;">
		<template #extra>
			<Button type="primary" @click="saveTemplate">保存</Button>
		</template>
	</page-header>
	<div style="height: 88vh;">
		<VueFlow :nodes="nodes" :edges="edges" :node-types="nodeTypes" :edge-types="edgeTypes"
			:connection-mode="ConnectionMode.Strict">
			<Background color="rgb(0,0,0)"/>
			<Panel :position="Position.Bottom">
				<Button type="primary" @click="openNewNodeModal">添加节点</Button>
			</Panel>
		</VueFlow>
	</div>

	<Modal :open="newNodeModalOpen" @cancel="closeNewNodeModal" @ok="newNodeConfirm">
		<template #title>选择节点</template>
		<Select :options="nodeTypeOptions" v-model:value="selectNodeType"></Select>
		<p>{{ nodeTypeOptions.find(n => n.value === selectNodeType).description }}</p>
		<template #okButton>
			<Button type="primary">确认</Button>
		</template>
	</Modal>

	<Drawer title="大模型配置" :open="llmDrawerOpen" @close="_=>{llmDrawerOpen = false;}">
    <LlmSetting v-model:node="currentSettingNode" :ref-options="settingRefOptions"
                :llm-list="llmList" :llm-options="llmOptions"/>
	</Drawer>
	<Drawer title="知识库检索配置" :open="knowledgeRetrievalDrawerOpen" @close="_=>{knowledgeRetrievalDrawerOpen = false;}"></Drawer>
	<Drawer title="知识库写入配置" :open="knowledgeWriteDrawerOpen" @close="_=>{knowledgeWriteDrawerOpen = false;}"></Drawer>
  <Drawer title="开始配置" :open="startDrawerOpen" @close="_=>{startDrawerOpen = false;}">
    <start-setting v-model:node="currentSettingNode"/>
  </Drawer>
</template>

<style>
@import '@vue-flow/core/dist/style.css';

</style>