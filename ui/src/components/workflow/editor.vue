<script setup>
import {ref, watch} from 'vue'
import { Panel, ConnectionMode, Position, useVueFlow, VueFlow, MarkerType } from '@vue-flow/core'
import { Background } from '@vue-flow/background';
import {Button, PageHeader, Select, Modal, Drawer, Input, message} from 'ant-design-vue';
import llmAPI from '../../api/llm';
import templateAPI from '../../api/template.js';
import workflowAPI from '../../api/workflow.js';
import LlmSetting from './setting/LLMSetting.vue';
import StartSetting from "./setting/StartSetting.vue";
import NodeUtil from "../../util/nodeUtil.js";
import EndSetting from "./setting/EndSetting.vue";
import ExecutionLog from "../instance/executionLog.vue";
import {useRoute, useRouter} from "vue-router";
import types from "./types.js";
import CrawlerSetting from "./setting/CrawlerSetting.vue";

const props = defineProps(['isNewTemplate','template'])
const route = useRoute();
const router = useRouter();
const nodeTypeOptions = types.nodeTypeOptions;
const edgeTypes = types.edgeTypes;

const selectNodeType = ref("llm");

const nodeTypes = types.nodeTypes;

const nodes = ref([]);
const edges = ref([]);

if (props.isNewTemplate) {
  props.template.name = "新建模板"
  nodes.value = [
    {
      id: "start",
      type: "start",
      position: { x: 100, y: 100 },
      name: "开始",
      data: {
        startNodeData: {
          inputVariables: [
            {name: "input", type: "string", "value": ""}
          ]
        }
      }
    },
    {
      id: "end",
      type: "end",
      name: "结束",
      position: { x: 200, y: 100 },
      data: {
        endNodeData: {
          outputVariables: [
            {name: "output", type: "string", value: ""}
          ]
        },
      }
    }
  ];
}else if (props.template['data']) {
  const definition = JSON.parse(props.template['data']);
  nodes.value = definition.nodes;
  edges.value = definition.edges;
}
watch(()=>props.template, function (oldVal, newVal) {
  const definition = JSON.parse(props.template['data']);
  nodes.value = definition.nodes;
  edges.value = definition.edges;
}, {deep: true});

const newNodeModalOpen = ref(false);
const { onNodeClick, onConnect, onEdgesChange, onNodeDragStop, 
	onNodesChange} = useVueFlow();

const executeLogDrawerOpen = ref(false);
const llmDrawerOpen = ref(false);
const knowledgeRetrievalDrawerOpen = ref(false);
const knowledgeWriteDrawerOpen = ref(false);
const startDrawerOpen = ref(false);
const endDrawerOpen = ref(false);
const crawlerDrawerOpen = ref(false);
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
    case "crawler": crawlerDrawerOpen.value = true; break;
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
	types.createNodeData(node);
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
    name: props.template.name,
    data: data,
  }).then(resp=>{
    router.push("/editor/"+resp.data['id']);
  })
}

function updateTemplate() {
  const data = getJSON();
  props.template.data = data;
  templateAPI.updateTemplate(props.template).then(resp=>{
    message.success("更新成功");
  }).catch(_=>{
    message.error("更新失败");
  })
}
// 获取前驱节点的输出变量列表
function getPrevNodesOutputs() {
  const prevNodes = NodeUtil.getPrevNodes(currentSettingNode.value.id, nodes.value, edges.value);
  let options = [];
  prevNodes.forEach(node=>{
    if (!node.data) return;
    let outputVariables;
    switch (node.type) {
      case "llm": outputVariables = node.data['llmNodeData'].outputVariables; break;
      case "start": outputVariables = node.data['startNodeData'].inputVariables; break;
      case "end": outputVariables = node.data['endNodeData'].outputVariables; break;
      case "crawler": outputVariables = node.data['crawlerNodeData'].outputVariables; break;
    }
    if (outputVariables) {
      let option = {
        label: node.name,
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
    llmList.value = resp.data;
    const options = [];
    llmList.value.forEach(item=>{
      options.push({label: item.name, value: item.id});
    });
    llmOptions.value = options;
    llmDrawerOpen.value = true;
  });
}

const outputInterval = ref(0);
const executeOutputs = ref([]);
// 创建临时实例运行
function execute() {
  const request = {
    inputs: {
      input: ""
    }
  };
  if (!props.isNewTemplate) {
    request['templateId'] = props.template.id;
  }else {
    request['definition'] = getJSON();
  }
  workflowAPI.start(request).then(resp=>{
    executeLogDrawerOpen.value = true;
    outputInterval.value = setInterval(_=>getExecuteOutputs(resp.data['workflowId']), 1000);
  });
}

function getExecuteOutputs(workflowId) {
  workflowAPI.getOutputs(workflowId).then(resp=>{
    executeOutputs.value = resp.data;
    const endNode = resp.data.find(n=>n.type === 'end');
    if (endNode) {
      clearInterval(outputInterval.value);
    }
  });
}
</script>

<template>
	<page-header title="流程设计" style="border: 1px solid rgb(235, 237, 240); height: 10vh;">
		<template #extra>
      <Input v-model:value="template.name"></Input>
			<Button type="primary" @click="saveTemplate" v-if="isNewTemplate">保存</Button>
      <Button type="primary" v-else @click="updateTemplate">更新</Button>
      <Button type="primary" success @click="execute">运行</Button>
		</template>
	</page-header>
	<div style="height: 88vh;">
		<VueFlow :nodes="nodes" :edges="edges" :node-types="nodeTypes" :edge-types="edgeTypes"
			:connection-mode="ConnectionMode.Strict" :edges-updatable="true">
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
  <Drawer title="结果配置" :open="endDrawerOpen" @close="_=>{endDrawerOpen = false;}">
    <end-setting :output-variables="currentSettingNode.data['endNodeData'].outputVariables"
                 :ref-options="settingRefOptions" :node="currentSettingNode"/>
  </Drawer>
  <Drawer title="爬虫配置" :open="crawlerDrawerOpen" @close="_=>{crawlerDrawerOpen = false;}">
    <CrawlerSetting :ref-options="settingRefOptions" :node="currentSettingNode"/>
  </Drawer>
  <Drawer title="执行结果" :open="executeLogDrawerOpen" @close="_=>{executeLogDrawerOpen = false;}">
    <execution-log :outputs="executeOutputs"></execution-log>
  </Drawer>
</template>

<style>
@import '@vue-flow/core/dist/style.css';

</style>