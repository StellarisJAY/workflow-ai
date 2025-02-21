<script setup>
import {
  Button,
  Drawer,
  Form,
  FormItem,
  Input,
  message,
  PageHeader,
  Spin,
  Table,
  Tag,
  UploadDragger,
    InputNumber
} from "ant-design-vue";
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

const documentListLoading = ref(false);

const processOptionDrawerOpen = ref(false);
const processOptions = ref({});


onMounted(()=>{
  listDocuments();
});

function listDocuments() {
  knowledgeBaseAPI.listDocuments(kbId, query.value).then((resp) => {
    if (resp.data && resp.data.length > 0) {
      kbDocumentList.value = resp.data;
    }else {
      kbDocumentList.value = [];
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

function deleteDocument(id) {
  documentListLoading.value = true;
  knowledgeBaseAPI.deleteFile(id).then((resp) => {
    message.success("删除成功");
    listDocuments();
    documentListLoading.value = false;
  }).catch((err) => {
    message.error("删除失败");
    documentListLoading.value = false;
  });
}

function openProcessOptionsDrawer(id) {
  knowledgeBaseAPI.getFileProcessOptions(id).then((resp) => {
    processOptions.value = resp.data;
    const separators = JSON.parse(processOptions.value['separators']);
    processOptions.value.separatorsConverted = escapeString(separators.join(""));
    processOptionDrawerOpen.value = true;
  });
}

function updateProcessOptions() {
  const unescaped = unescapeString(processOptions.value.separatorsConverted);
  processOptions.value.separators = unescaped.split("");
  knowledgeBaseAPI.updateFileProcessOptions(processOptions.value).then((resp) => {
    message.success("更新成功");

  });
}

function escapeString(str) {
  return str.replace(/[\\]/g, '\\\\')
      .replace(/[\/]/g, '\\/')
      .replace(/[\b]/g, '\\b')
      .replace(/[\f]/g, '\\f')
      .replace(/[\n]/g, '\\n')
      .replace(/[\r]/g, '\\r')
      .replace(/[\t]/g, '\\t')
      .replace(/[\"]/g, '\\"')
      .replace(/[\']/g, "\\'");
}

function unescapeString(str) {
  return str.replace(/\\([\\\/bfnrt'"])/g, function(match, char) {
    switch (char) {
      case '\\': return '\\';
      case '/': return '/';
      case 'b': return '\b';
      case 'f': return '\f';
      case 'n': return '\n';
      case 'r': return '\r';
      case 't': return '\t';
      case '"': return '"';
      case "'": return "'";
      default: return char;
    }
  });
}

function startDocumentProcessing(id) {
  knowledgeBaseAPI.startFileProcessing(id).then((resp) => {
    message.success("开始解析");
  });
}
</script>

<template>
  <PageHeader>
    <template #extra>
      <Button type="primary" @click="openUploadDrawer">上传</Button>
    </template>
  </PageHeader>
  <Spin :spinning="documentListLoading">
    <Table :columns="columns" :data-source="kbDocumentList" :pagination="false">
      <template #bodyCell="{ column, _, record }">
        <template v-if="column.dataIndex === 'operation'">
          <a @click="startDocumentProcessing(record['id'])">解析</a>
          /
          <a @click="openProcessOptionsDrawer(record['id'])">设置</a>
          /
          <a @click="_=>{knowledgeBaseAPI.downloadFile(record['id'])}">下载</a>
          /
          <a @click="deleteDocument(record['id'])">删除</a>
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
  </Spin>
  <Drawer :open="uploadDrawerOpen" title="上传文档" @close="_=>{uploadDrawerOpen=false}">
    <Spin about="上传中" :spinning="uploading">
      <UploadDragger :multiple="false" v-model:file-list="uploadFileList" name="file"/>
      <Button type="primary" @click="uploadDocument">上传</Button>
    </Spin>
  </Drawer>
  <Drawer :open="processOptionDrawerOpen" title="解析设置" @close="_=>{processOptionDrawerOpen=false;}">
    <Form>
      <FormItem label="块token数">
        <InputNumber v-model:value="processOptions['chunkSize']" :min="128"/>
      </FormItem>
      <FormItem label="块分隔符">
        <Input v-model:value="processOptions['separatorsConverted']"/>
      </FormItem>
    </Form>
    <Button type="primary" @click="updateProcessOptions">修改</Button>
  </Drawer>
</template>

<style scoped>

</style>