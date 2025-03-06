<script setup>
import {useRoute, useRouter} from "vue-router";
import {onMounted, ref} from "vue";
import {Table, Pagination, Dropdown, Menu, MenuItem, Popconfirm, message,
  PageHeader, Button} from "ant-design-vue";
import workflowAPI from '../../api/workflow.js';
import TimeUtil from "../../util/timeUtil.js";
import {EditOutlined, PlayCircleOutlined, EllipsisOutlined} from "@ant-design/icons-vue";

const columns = [
  {title: "流程ID", dataIndex: "id", key: "id"},
  {title: "模板名称", dataIndex: "templateName", key: "templateName" },
  {title: "状态", dataIndex: "statusName", key: "statusName" },
  {title: "开始时间", dataIndex: "addTime", key: "addTime" },
  {title: "结束时间", dataIndex: "completeTime", key: "completeTime" },
  {title: "处理耗时", dataIndex: "duration", key: "duration" },
  {title: "操作", dataIndex: "operation", key: "operation" }
]
const router = useRouter();
const route = useRoute();
const workflowInstanceList = ref([]);
const query = ref({
  page: 1,
  pageSize: 10,
  paged: true,
});
const total = ref(0);

function listWorkflowInstance() {
  workflowAPI.list(query.value).then(resp=>{
    total.value = resp.total;
    workflowInstanceList.value = resp.data;
  });
}
function openWorkflowViewer(id) {
  router.push("/view/" + id);
}

onMounted(_=>{
  listWorkflowInstance();
});
</script>

<template>
  <PageHeader title="执行历史列表">
    <template #extra>
    </template>
  </PageHeader>

  <Table :pagination="false" :columns="columns" :data-source="workflowInstanceList" >
    <template #bodyCell="{ column, _, record }">
      <template v-if="column.dataIndex === 'operation'">
        <a @click="openWorkflowViewer(record.id)">详情</a>
      </template>
      <template v-if="column.dataIndex === 'addTime'">
        {{TimeUtil.formatDateTime(record.addTime)}}
      </template>
      <template v-if="column.dataIndex === 'completeTime'">
        {{TimeUtil.formatDateTime(record.completeTime)}}
      </template>
    </template>
  </Table>

  <Pagination v-model:current="query.page" :page-size="query.pageSize" :total="total" @change="listWorkflowInstance"/>
</template>

<style scoped>
</style>