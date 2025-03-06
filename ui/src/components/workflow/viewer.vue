<script setup>
import {useRoute, useRouter} from "vue-router";
import {onMounted, onUnmounted, ref} from "vue";
import {
  PageHeader, Button, Drawer, Form, FormItem, Collapse, CollapsePanel
} from "ant-design-vue";
import workflowAPI from '../../api/workflow.js';
import {ReloadOutlined} from "@ant-design/icons-vue";
import {useVueFlow, VueFlow} from "@vue-flow/core";
import {Background} from "@vue-flow/background";
import nodeConstants from './nodeConstants.js';
import TimeUtil from "../../util/timeUtil.js";
import NodeStatusTag from "./node/nodeStatusTag.vue";

const props = defineProps(['workflowId']);
const router = useRouter();
const route = useRoute();
const workflowInstance = ref({});
const nodeTypes = nodeConstants.nodeTypes;
const edgeTypes = nodeConstants.edgeTypes;

const nodes = ref([]);
const edges = ref([]);

const nodeDetailOpen = ref(false);
const executeLogOpen = ref(false);

const currentNode = ref({});
const currentNodeInstance = ref({});
const currentNodeOutput = ref({});
const currentNodeOutputVarTypes = ref({});

const {onNodeClick} = useVueFlow();

const queryInterval = ref(0);

onMounted(_=>{
  getWorkflowDetail();
});
onUnmounted(_=>{
  if (queryInterval.value !== 0) {
    clearInterval(queryInterval.value);
    queryInterval.value = 0;
  }
});
onNodeClick(nodeClickHandler);

function getWorkflowDetail() {
  workflowAPI.detail(props.workflowId).then(resp=>{
    workflowInstance.value = resp.data;
    const definition = JSON.parse(workflowInstance.value.data);
    nodes.value = definition.nodes;
    edges.value = definition.edges;
    const nodeStatusList = workflowInstance.value['nodeStatusList'];
    const passedEdges = workflowInstance.value['passedEdgesList'];
    const successBranches = workflowInstance.value['successBranchList'];
    // 设置节点状态名
    nodes.value.forEach(node=>{
      const n = nodeStatusList.find(ns=>ns.nodeId === node.id);
      if (n) {
        node.status = {id: n.status, text: n['statusName']};
      }else {
        node.status = {id: 0, text: "未到达"};
      }
      // 标记条件节点通过的分支
      if (node['type'] === 'condition') {
        const conditionNodeData = node['data']['conditionNodeData'];
        const branch = conditionNodeData.branches.find(b=>{
          return successBranches.findIndex(s=>s.nodeId === node['id'] && s['branch']===b['handle']) !== -1;
        });
        if (branch) {
          branch['success'] = true;
        }
      }
    });
    // 标记通过的连线
    edges.value.forEach(edge=>{
      const f = passedEdges.find(e=>e === edge.id);
      if (f) {
        edge.passed = true;
      }
    });
    if (queryInterval.value === 0 && workflowInstance.value.status === 0) {
      queryInterval.value = setInterval(_=>getWorkflowDetail(), 3000);
    } else if (queryInterval.value !== 0 && workflowInstance.value.status !== 0) {
      clearInterval(queryInterval.value);
      queryInterval.value = 0;
    }
  });
}

function openExecutionLog() {
  executeLogOpen.value = true;
}

function nodeClickHandler(event) {
  const node = event.node;
  workflowAPI.getNodeInstanceDetail(props.workflowId, node.id).then(resp=>{
    currentNode.value = node;
    currentNodeInstance.value = resp.data;
    if (currentNodeInstance.value.output) {
      currentNodeOutput.value = JSON.parse(currentNodeInstance.value.output);
    } else {
      currentNodeOutput.value = "";
    }
    currentNodeOutputVarTypes.value = currentNodeInstance.value['outputVariableTypes'];
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
  </Drawer>

  <Drawer size="large" title="节点详情" :open="nodeDetailOpen" @close="_=>{nodeDetailOpen = false;}">
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
    <Collapse>
      <CollapsePanel v-for="(value, key) in currentNodeOutput" :key="key" :header="key">
        <p v-if="currentNodeOutputVarTypes[key]==='string' || currentNodeOutputVarTypes[key]==='number'">{{value}}</p>
        <Collapse v-if="currentNodeOutputVarTypes[key]==='array_str' || currentNodeOutputVarTypes[key]==='array_num'">
          <CollapsePanel v-for="(item, i) in value" :header="i">
            <p>{{item}}</p>
          </CollapsePanel>
        </Collapse>
      </CollapsePanel>
    </Collapse>
    <h4 v-if="currentNodeInstance['error']">错误信息</h4>
    <p v-if="currentNodeInstance['error']">{{currentNodeInstance['error']}}</p>
  </Drawer>
</template>

<style scoped>
</style>