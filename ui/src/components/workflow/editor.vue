<script setup>
import { markRaw, ref } from 'vue'
import { Panel, ConnectionMode, Position, useVueFlow, VueFlow, MarkerType } from '@vue-flow/core'
import { Background } from '@vue-flow/background';
import CustomNode from './customNode.vue';
import LLMNode from './llmNode.vue';
import CustomEdge from './customEdge.vue';
import StartNode from './startNode.vue';
import EndNode from './endNode.vue';
import knowledgeRetrivalNode from './knowledgeRetrivalNode.vue';
import knowledgeWriteNode from './knowledgeWriteNode.vue';
import conditionNode from './conditionNode.vue';
import { Button, PageHeader, Select, Modal, Drawer } from 'ant-design-vue';

const nodeTypeOptions = ref([
	{ value: "llm", label: "大模型", description: "使用提示词和变量让大模型生成内容" },
	{ value: "knowledgeRetrival", label: "知识库检索", description: "从知识库检索数据" },
	{ value: "knowledgeWrite", label: "知识库写入", description: "将数据写入知识库" },
	{ value: "condition", label: "条件", description: "条件分支" }
]);

const selectNodeType = ref("llm");

const nodeTypes = ref({
	custom: markRaw(CustomNode),
	llm: markRaw(LLMNode),
	start: markRaw(StartNode),
	end: markRaw(EndNode),
	condition: markRaw(conditionNode),
	knowledgeRetrival: markRaw(knowledgeRetrivalNode),
	knowledgeWrite: markRaw(knowledgeWriteNode)
});

const edgeTypes = ref({
	custom: markRaw(CustomEdge),
});

const nodes = ref([
	{
		id: "1",
		type: "start",
		position: { x: 100, y: 100 },
		data: {}
	},
	{
		id: "2",
		type: "end",
		position: { x: 200, y: 100 },
		data: {}
	}
]);
const edges = ref([]);
const newNodeModalOpen = ref(false);
const { onNodeClick, onConnect, onEdgesChange, onNodeDragStop } = useVueFlow();

const llmDrawerOpen = ref(false);
const knowledgeRetrivalDrawerOpen = ref(false);
const knowledgeWriteDrawerOpen = ref(false);
const currentSettingNode = ref({});

onNodeClick(event => {
	currentSettingNode.value = event.node;
	switch (event.node.type) {
		case "llm": llmDrawerOpen.value = true; break;
		case "knowledgeRetrival": knowledgeRetrivalDrawerOpen.value = true; break;
		case "knowledgeWrite": knowledgeWriteDrawerOpen.value = true; break;
	}
});
onConnect(event => {
	edges.value.push({
		id: crypto.randomUUID(),
		type: "custom",
		source: event.source,
		target: event.target,
		markerStart: MarkerType.ArrowClosed,
	});
});

onEdgesChange(ev => {
	ev.forEach(e => {
		if (e.type === "remove") {
			removeNode(e.id);
		}
	});
});

onNodeDragStop(ev => {
	nodes.value.find(n => n.id === ev.node.id).position = ev.node.position;
})

function addNode(nodeType) {
	const id = crypto.randomUUID();
	nodes.value.push({
		type: nodeType,
		id: id,
		position: { x: 300, y: 200 },
		data: {}
	});
}

function removeNode(id) {
	edges.value = edges.value.filter(n => n.id !== id);
}

function getJSON() {
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

function saveTemplate() {
	const data = getJSON();
	console.log(data);
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
			<Background />
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
	</Drawer>
	<Drawer title="知识库检索配置" :open="knowledgeRetrivalDrawerOpen" @close="_=>{knowledgeRetrivalDrawerOpen = false;}"></Drawer>
	<Drawer title="知识库写入配置" :open="knowledgeWriteDrawerOpen" @close="_=>{knowledgeWriteDrawerOpen = false;}"></Drawer>
</template>

<style>
@import '@vue-flow/core/dist/style.css';
@import '@vue-flow/core/dist/theme-default.css';
</style>