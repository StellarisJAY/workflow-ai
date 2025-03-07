<script setup>
import {ref, watch} from 'vue'
import {ConnectionMode, MarkerType, Panel, Position, useVueFlow, VueFlow} from '@vue-flow/core'
import {Background} from '@vue-flow/background';
import {Button, Drawer, Input, message, Modal, PageHeader, Select} from 'ant-design-vue';
import templateAPI from '../../api/template.js';
import LlmSetting from './setting/LLMSetting.vue';
import StartSetting from "./setting/StartSetting.vue";
import EndSetting from "./setting/EndSetting.vue";
import {useRoute, useRouter} from "vue-router";
import nodeConstants from "./nodeConstants.js";
import CrawlerSetting from "./setting/CrawlerSetting.vue";
import ConditionSetting from "./setting/ConditionSetting.vue";
import {randomUUID} from "../../util/uuid.js";
import KnowledgeRetrievalSetting from "./setting/KnowledgeRetrievalSetting.vue";
import NodeConstants from "./nodeConstants.js";
import WebSearchSetting from "./setting/WebSearchSetting.vue";
import KeywordExtractionSetting from "./setting/KeywordExtractionSetting.vue";
import QuestionOptimizationSetting from "./setting/QuestionOptimizationSetting.vue";
import ImageUnderstandingSetting from "./setting/ImageUnderstandingSetting.vue";
import OCRSetting from "./setting/OCRSetting.vue";

const props = defineProps(['isNewTemplate','template'])
const route = useRoute();
const router = useRouter();
const nodeTypeOptions = nodeConstants.nodeTypeOptions;
const edgeTypes = nodeConstants.edgeTypes;
const { onNodeClick, onConnect, onEdgesChange, onNodeDragStop,
  onNodesChange} = useVueFlow();

const selectNodeType = ref("llm");

const nodeTypes = nodeConstants.nodeTypes;

const nodes = ref([]);
const edges = ref([]);

if (props.isNewTemplate) {
  props.template.name = "新建模板"
  nodes.value = [
    {
      id: "start",
      type: "start",
      position: { x: 100, y: 100 },
      data: {
        name: "开始",
        defaultAllowVarTypes: ["string", "number", "image_file"],
        allowAddInputVar: true,
        startNodeData: {
          inputVariables: [
            {name: "input", type: "string", value: "", allowRef: false, isRef: false}
          ]
        }
      }
    },
    {
      id: "end",
      type: "end",
      defaultAllowVarTypes: ["string", "number"],
      position: { x: 200, y: 100 },
      data: {
        name: "结束",
        endNodeData: {
          outputVariables: [
          ]
        },
      }
    }
  ];
} else if (props.template['data']) {
  const definition = JSON.parse(props.template['data']);
  nodes.value = definition.nodes;
  edges.value = definition.edges;
}
watch(()=>props.template, function () {
  const definition = JSON.parse(props.template['data']);
  nodes.value = definition.nodes;
  edges.value = definition.edges;
}, {deep: true});

const newNodeModalOpen = ref(false);

const llmDrawerOpen = ref(false);
const knowledgeRetrievalDrawerOpen = ref(false);
const knowledgeWriteDrawerOpen = ref(false);
const startDrawerOpen = ref(false);
const endDrawerOpen = ref(false);
const crawlerDrawerOpen = ref(false);
const conditionDrawerOpen = ref(false);
const webSearchDrawerOpen = ref(false);
const keywordExtractionDrawerOpen = ref(false);
const questionOptimizationDrawerOpen = ref(false);
const imageUnderstandingDrawerOpen = ref(false);
const ocrDrawerOpen = ref(false);

const currentSettingNodes = ref({});

// 点击节点，弹出侧边设置
onNodeClick(event => {
  const curNode = event.node;
  currentSettingNodes.value[curNode.type] = curNode;
	switch (curNode.type) {
		case "llm": llmDrawerOpen.value=true; break;
		case "knowledgeRetrieval": knowledgeRetrievalDrawerOpen.value=true; break;
		case "knowledgeWrite": knowledgeWriteDrawerOpen.value = true; break;
    case "start": startDrawerOpen.value = true; break;
    case "end": endDrawerOpen.value = true; break;
    case "crawler": crawlerDrawerOpen.value = true; break;
    case "condition": conditionDrawerOpen.value = true; break;
    case "webSearch": webSearchDrawerOpen.value = true; break;
    case "keywordExtraction": keywordExtractionDrawerOpen.value = true; break;
    case "questionOptimization": questionOptimizationDrawerOpen.value = true; break;
    case "imageUnderstanding": imageUnderstandingDrawerOpen.value = true; break;
    case "ocr": ocrDrawerOpen.value = true; break;
	}
});
// 节点连线事件，添加edge
onConnect(event => {
  const exist = edges.value.find(e=>{return e['source'] === event['source'] &&
      e['target'] === event['target'] &&
      e['sourceHandle'] === event['sourceHandle'] &&
      e['targetHandle'] === event['targetHandle']});
  if (exist) {
    return;
  }
	edges.value.push({
		id: randomUUID(),
		type: "custom",
		source: event.source,
		target: event.target,
		markerEnd: MarkerType.ArrowClosed,
    sourceHandle: event.sourceHandle,
    targetHandle: event.targetHandle,
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
	const id = randomUUID();
	const node = {
		type: nodeType,
		id: id,
		position: { x: 300, y: 200 },
		data: {}
	};
	nodeConstants.createNodeData(node);
	nodes.value.push(node);
}
// 删除节点
function removeNode(id) {
	nodes.value = nodes.value.filter(n => n.id !== id);
  dispatchEvent(new CustomEvent(NodeConstants.deleteNodeEvent, {
    detail: {
      nodeId: id,
    }
  }));
}
// 删除连线
function removeEdge(id) {
	edges.value = edges.value.filter(e=>e.id !== id);
  dispatchEvent(new CustomEvent(NodeConstants.deleteEdgeEvent, {
    detail: {
      edgeId: id,
    }
  }));
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
  nodeConstants.loadNodePrototypes().then(()=>{
    newNodeModalOpen.value = true;
  });
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
  props.template.data = getJSON();
  templateAPI.updateTemplate(props.template).then(resp=>{
    message.success("更新成功");
  }).catch(_=>{
    message.error("更新失败");
  })
}
</script>

<template>
	<page-header title="流程设计" style="border: 1px solid rgb(235, 237, 240); height: 10vh;">
		<template #extra>
      <Input v-model:value="template.name"></Input>
			<Button type="primary" @click="saveTemplate" v-if="isNewTemplate">保存</Button>
      <Button type="primary" v-else @click="updateTemplate">更新</Button>
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
		<Select :options="nodeTypeOptions" v-model:value="selectNodeType" style="width: 200px;"></Select>
		<p>{{ nodeTypeOptions.find(n => n.value === selectNodeType).description }}</p>
		<template #okButton>
			<Button type="primary">确认</Button>
		</template>
	</Modal>

	<Drawer title="大模型配置" size="large"
          :open="llmDrawerOpen"
          @close="_=>{llmDrawerOpen = false;}"
          :destroy-on-close="true">
    <LlmSetting :node="currentSettingNodes['llm']"/>
	</Drawer>
	<Drawer title="知识库检索配置" size="large"
          :open="knowledgeRetrievalDrawerOpen"
          @close="_=>{knowledgeRetrievalDrawerOpen = false;}"
          :destroy-on-close="true">
    <KnowledgeRetrievalSetting :node="currentSettingNodes['knowledgeRetrieval']"/>
  </Drawer>
	<Drawer title="知识库写入配置" size="large"
          :open="knowledgeWriteDrawerOpen"
          @close="_=>{knowledgeWriteDrawerOpen = false;}"
          :destroy-on-close="true"></Drawer>
  <Drawer title="开始配置" size="large"
          :open="startDrawerOpen"
          @close="_=>{startDrawerOpen = false;}"
          :destroy-on-close="true">
    <start-setting v-model:node="currentSettingNodes['start']"/>
  </Drawer>
  <Drawer title="结果配置" size="large"
          :open="endDrawerOpen"
          @close="_=>{endDrawerOpen = false;}"
          :destroy-on-close="true">
    <end-setting :output-variables="currentSettingNodes['end'].data['endNodeData'].outputVariables"
                 :node="currentSettingNodes['end']"/>
  </Drawer>
  <Drawer title="爬虫配置" size="large"
          :open="crawlerDrawerOpen"
          @close="_=>{crawlerDrawerOpen = false;}"
          :destroy-on-close="true">
    <CrawlerSetting :node="currentSettingNodes['crawler']"/>
  </Drawer>
  <Drawer title="条件设置" size="large"
          :open="conditionDrawerOpen"
          @close="_=>{conditionDrawerOpen = false;}"
          :destroy-on-close="true">
    <ConditionSetting :node="currentSettingNodes['condition']"/>
  </Drawer>
  <Drawer title="搜索设置" size="large"
          :open="webSearchDrawerOpen"
          @close="_=>{webSearchDrawerOpen = false;}"
          :destroy-on-close="true">
    <WebSearchSetting :node="currentSettingNodes['webSearch']"/>
  </Drawer>
  <Drawer title="关键词提取设置" size="large"
          :open="keywordExtractionDrawerOpen"
          @close="_=>{keywordExtractionDrawerOpen = false;}"
          :destroy-on-close="true">
    <keywordExtractionSetting :node="currentSettingNodes['keywordExtraction']"/>
  </Drawer>
  <Drawer title="提问优化设置" size="large"
          :open="questionOptimizationDrawerOpen"
          @close="_=>{questionOptimizationDrawerOpen = false;}"
          :destroy-on-close="true">
    <QuestionOptimizationSetting :node="currentSettingNodes['questionOptimization']"/>
  </Drawer>

  <Drawer title="图像理解设置" size="large"
          :open="imageUnderstandingDrawerOpen"
          @close="_=>{imageUnderstandingDrawerOpen = false;}"
          :destroy-on-close="true">
    <ImageUnderstandingSetting :node="currentSettingNodes['imageUnderstanding']"/>
  </Drawer>

  <Drawer title="文字提取设置" size="large"
          :open="ocrDrawerOpen"
          @close="_=>{ocrDrawerOpen = false;}"
          :destroy-on-close="true">
    <OCRSetting :node="currentSettingNodes['ocr']"/>
  </Drawer>
</template>

<style>
@import '@vue-flow/core/dist/style.css';

</style>