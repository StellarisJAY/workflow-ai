<script setup>
import {useRoute, useRouter} from "vue-router";
import {onMounted, ref} from "vue";
import {
  Table, Pagination, Dropdown, Menu, MenuItem, Popconfirm, message,
  PageHeader, Button, Input, Drawer, Form, FormItem
} from "ant-design-vue";
import workflowAPI from '../../api/workflow.js';
import {ReloadOutlined} from "@ant-design/icons-vue";
import {useVueFlow, VueFlow} from "@vue-flow/core";
import {Background} from "@vue-flow/background";
import types from './types.js';
import ExecutionLog from "../instance/executionLog.vue";
import TimeUtil from "../../util/timeUtil.js";
import NodeStatusTag from "./node/nodeStatusTag.vue";

const props = defineProps(['workflowId']);
const router = useRouter();
const route = useRoute();
const workflowInstance = ref({});
const nodeTypes = types.nodeTypes;
const edgeTypes = types.edgeTypes;

const nodes = ref([]);
const edges = ref([]);

const nodeDetailOpen = ref(false);
const executeLogOpen = ref(false);
const instanceOutputs = ref([]);

const currentNode = ref({});
const currentNodeInstance = ref({});
const currentNodeOutput = ref({});

const {onNodeClick} = useVueFlow();

onMounted(_=>{
  getWorkflowDetail();
  getInstanceOutputs();
});
onNodeClick(nodeClickHandler);

function getWorkflowDetail() {
  workflowAPI.detail(props.workflowId).then(resp=>{
    workflowInstance.value = resp.data;
    const definition = JSON.parse(workflowInstance.value.data);
    nodes.value = definition.nodes;
    edges.value = definition.edges;
    const nodeStatusList = workflowInstance.value.nodeStatusList;
    nodes.value.forEach(node=>{
      const n = nodeStatusList.find(ns=>ns.nodeId === node.id);
      if (n) {
        node.status = {id: n.status, text: n.statusName};
      }else {
        node.status = {id: 0, text: "未到达"};
      }
    });
  });
}

function getInstanceOutputs() {
  workflowAPI.getOutputs(props.workflowId).then(resp=>{
    instanceOutputs.value = resp.data;
  });
}

function openExecutionLog() {
  getInstanceOutputs();
  executeLogOpen.value = true;
}

function nodeClickHandler(event) {
  const node = event.node;
  workflowAPI.getNodeInstanceDetail(props.workflowId, node.id).then(resp=>{
    currentNode.value = node;
    currentNodeInstance.value = resp.data;
    currentNodeOutput.value = JSON.parse(currentNodeInstance.value.output);
    nodeDetailOpen.value = true;
  });
}
</script>

<template>
  <page-header title="流程详情" style="border: 1px solid rgb(235, 237, 240); height: 10vh;">
    <template #extra>
      <Button @click="getWorkflowDetail">
        <ReloadOutlined />
      </Button>
      <Button type="primary" @click="openExecutionLog">日志</Button>
    </template>
  </page-header>
  <div style="height: 88vh;">
    <VueFlow :nodes="nodes" :edges="edges" :node-types="nodeTypes" :edge-types="edgeTypes"
             :edges-updatable="false"
             :nodes-connectable="false"
             :nodes-draggable="false">
      <Background color="rgb(0,0,0)"/>
    </VueFlow>
  </div>

  <Drawer title="执行日志" :open="executeLogOpen" @close="_=>{executeLogOpen = false;}">
    <execution-log :outputs="instanceOutputs"/>
  </Drawer>

  <Drawer title="节点详情" :open="nodeDetailOpen" @close="_=>{nodeDetailOpen = false;}">
    <Form>
      <FormItem label="状态">
        <node-status-tag :status="{id: currentNodeInstance.status, text: currentNodeInstance.statusName}"/>
      </FormItem>
      <FormItem label="开始时间">{{TimeUtil.formatDateTime(currentNodeInstance['addTime'])}}</FormItem>
      <FormItem label="结束时间">
        {{
          ( currentNodeInstance.status === 0||currentNodeInstance.status === 1 )?'/':TimeUtil.formatDateTime(currentNodeInstance['completeTime'])
        }}
      </FormItem>
    </Form>
    <h4 v-if="currentNodeOutput">输出变量</h4>
    <Form>
      <FormItem v-for="(value, key) in currentNodeOutput" :label="key">
        {{value}}
      </FormItem>
    </Form>
    <h4 v-if="currentNodeInstance['error']">错误信息</h4>
    <p v-if="currentNodeInstance['error']">{{currentNodeinstance['error']}}</p>
  </Drawer>
</template>

<style scoped>
</style>