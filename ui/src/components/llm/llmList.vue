<script setup>
import {useRoute, useRouter} from "vue-router";
import {onMounted, ref} from "vue";
import {Table, Pagination, Form, FormItem,message, InputPassword,
  PageHeader, Button, Drawer, Input, Select} from "ant-design-vue";
import llmAPI from '../../api/llm.js';
import {llmTypeOptions, llmAPITypeOptions} from "../../api/const.js";

const columns = [
  {title: "名称", dataIndex: "name", key: "name" },
  {title: "模型", dataIndex: "code", key: "code" },
  {title: "类型", dataIndex: "modelType", key: "modelType" },
  {title: "API类型", dataIndex: "apiType", key: "apiType" },
  {title: "操作", dataIndex: "operation", key: "operation" }
]
const router = useRouter();
const route = useRoute();
const llmList = ref([]);
const query = ref({
  page: 1,
  pageSize: 10,
});
const total = ref(0);

const newLLM = ref({});
const addLLMDrawerOpen = ref(false);

function listLLM() {
  llmAPI.listModels(query).then(resp=>{
    total.value = resp.data.total;
    llmList.value = resp.data;
  });
}

function openAddLLMDrawer() {
  newLLM.value = {
    name: "",
    modelType: "chat",
    apiType: "openai",
    baseUrl: "http://example.com/",
    apiKey: "",
    code: "",
  };
  addLLMDrawerOpen.value = true;
}

function createLLM() {
  llmAPI.createModel(newLLM.value).then(_=>{
    listLLM();
    message.success("创建成功");
    addLLMDrawerOpen.value = false;
  }).catch(err=>{
    message.error("创建失败");
  })
}

onMounted(_=>{
  listLLM();
});
</script>

<template>
  <PageHeader title="大模型列表">
    <template #extra>
      <Button type="primary" @click="openAddLLMDrawer">添加</Button>
    </template>
  </PageHeader>

  <Table :pagination="false" :columns="columns" :data-source="llmList" >
    <template #bodyCell="{ column, _, record }">
      <template v-if="column.dataIndex === 'operation'">
        <a>详情</a>
      </template>
    </template>
  </Table>

  <Pagination :current="query.page" :page-size="query.pageSize" :total="total"></Pagination>

  <Drawer title="添加大模型" :open="addLLMDrawerOpen" @close="_=>{addLLMDrawerOpen=false;}">
    <Form>
      <FormItem label="名称">
        <Input v-model:value="newLLM.name"/>
      </FormItem>
      <FormItem label="类型">
        <Select v-model:value="newLLM.modelType" :options="llmTypeOptions"/>
      </FormItem>
      <FormItem label="模型">
        <Input v-model:value="newLLM.code"/>
      </FormItem>
      <FormItem label="API类型">
        <Select v-model:value="newLLM.apiType" :options="llmAPITypeOptions"/>
      </FormItem>
      <FormItem label="API URL">
        <Input v-model:value="newLLM.baseUrl"/>
      </FormItem>
      <FormItem label="API Key">
        <InputPassword v-model:value="newLLM.apiKey"/>
      </FormItem>
      <Button type="primary" @click="createLLM">创建</Button>
    </Form>
  </Drawer>
</template>

<style scoped>
</style>