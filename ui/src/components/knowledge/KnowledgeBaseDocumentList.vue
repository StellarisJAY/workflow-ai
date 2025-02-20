<script setup>
import {Table, PageHeader, Button, UploadDragger, Drawer, message, Spin, Tag} from "ant-design-vue";
import {useRoute} from "vue-router";
import {onMounted, ref} from "vue";
import knowledgeBaseAPI from "../../api/knowledgeBase.js";
import TimeUtil from "../../util/timeUtil.js";

const route = useRoute();
const kbId = route.params['id'];

const kbDocumentList = ref([]);
const query = ref({
  page: 1,
  pageSize: 10,
  paged: true
});
const total = ref(0);
const columns = [
  {title: "文件名", dataIndex: "name", key: "name" },
  {title: "大小", dataIndex: "length", key: "length" },
  {title: "上传时间", dataIndex: "addTime", key: "addTime" },
  {title: "状态", dataIndex: "status", key: "status" },
  {title: "操作", dataIndex: "operation", key: "operation" },
];

const uploadDrawerOpen = ref(false);
const uploadFileList = ref([]);
const uploading = ref(false);

onMounted(()=>{
  listDocuments();
});

function listDocuments() {
  knowledgeBaseAPI.listDocuments(kbId, query.value).then((resp) => {
    if (resp.data && resp.data.length > 0) {
      kbDocumentList.value = resp.data;
    }
    total.value = resp.data.length;
  });
}

function uploadDocument() {
  const formData = new FormData();
  formData.append("file", uploadFileList.value[0].originFileObj);
  formData.append("kbId", kbId);
  uploading.value = true;
  knowledgeBaseAPI.upload(formData).then((resp) => {
    message.success("上传成功");
    uploadDrawerOpen.value = false;
    uploading.value = false;
    listDocuments();
  }).catch((err) => {
    message.error(err);
    uploading.value = false;
  })
}

function openUploadDrawer() {
  uploadFileList.value = [];
  uploading.value = false;
  uploadDrawerOpen.value = true;
}
</script>

<template>
  <PageHeader>
    <template #extra>
      <Button type="primary" @click="openUploadDrawer">上传</Button>
    </template>
  </PageHeader>
  <Table :columns="columns" :data-source="kbDocumentList" :pagination="false">
    <template #bodyCell="{ column, _, record }">
      <template v-if="column.dataIndex === 'operation'">
        <a>解析</a>
        /
        <a>详情</a>
        /
        <a>删除</a>
      </template>
      <template v-if="column.dataIndex === 'addTime'">
        {{TimeUtil.formatDateTime(record['addTime'])}}
      </template>
      <template v-if="column.dataIndex === 'status'">
        <Tag v-if="record['status'] === 0" color="red">{{record['statusName']}}</Tag>
        <Tag v-else-if="record['status'] === 1" color="yellow">{{record['statusName']}}</Tag>
        <Tag v-else-if="record['status'] === 2" color="green">{{record['statusName']}}</Tag>
        <Tag v-else-if="record['status'] === 3" color="red">{{record['statusName']}}</Tag>
      </template>
    </template>
  </Table>
  <Drawer :open="uploadDrawerOpen" title="上传文档" @close="_=>{uploadDrawerOpen=false}">
    <Spin about="上传中" :spinning="uploading">
      <UploadDragger :multiple="false" v-model:file-list="uploadFileList" name="file"/>
      <Button type="primary" @click="uploadDocument">上传</Button>
    </Spin>
  </Drawer>
</template>

<style scoped>

</style>