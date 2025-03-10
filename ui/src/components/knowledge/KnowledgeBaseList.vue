<script setup>
import {
  Card,
  List,
  ListItem,
  Button,
  Tag,
  Drawer,
  Popconfirm,
  Form,
  PageHeader,
  Dropdown,
  Pagination, FormItem, Menu, Input, Textarea, Select, message,
} from "ant-design-vue";
import {onMounted, ref} from "vue";
import llmAPI from '../../api/llm.js';
import knowledgeBaseAPI from '../../api/knowledgeBase.js';
import {SettingOutlined, UploadOutlined} from "@ant-design/icons-vue";
import timeUtil from "../../util/timeUtil.js";
import {useRouter} from "vue-router";
const query = ref({
  page: 1,
  pageSize: 16,
});
const total = ref(0);
const knowledgeBaseList = ref([]);

const newKnowledgeBase = ref({});
const createKBDrawerOpen = ref(false);
const embeddingModels = ref([]);
const router = useRouter();
onMounted(_=>{
  listKnowledgeBase();
});

function openCreateKBDrawer() {
  listEmbeddingModels();
  newKnowledgeBase.value = {name: "", description: ""};
  createKBDrawerOpen.value = true;
}

function listEmbeddingModels() {
  llmAPI.listModels({modelType: "embedding"}).then(resp=>{
    const models = resp.data;
    const options = [];
    models.forEach(model=>{
      options.push({label: model.name, value: model.id});
    });
    embeddingModels.value = options;
  });
}

function listKnowledgeBase() {
  knowledgeBaseAPI.list(query.value).then(resp=>{
    if (resp.data && resp.data.length > 0) {
      knowledgeBaseList.value = resp.data;
    }
    total.value = resp.total;
  });
}

function createKnowledgeBase() {
  knowledgeBaseAPI.create(newKnowledgeBase.value).then(resp=>{
    message.success("创建成功");
  });
}
</script>

<template>
  <PageHeader title="知识库">
    <template #extra>
      <Button type="primary" @click="openCreateKBDrawer">创建</Button>
    </template>
  </PageHeader>
  <List :grid="{ gutter: 4, column: 4 }" :data-source="knowledgeBaseList">
    <template #renderItem="{item}">
      <ListItem>
        <Card :title="item.name" :hoverable="true">
          <template #actions>
            <UploadOutlined />
            <SettingOutlined @click="router.push('/knowledgeBase/'+item.id)"/>
          </template>
          <template #extra>
          </template>
          <FormItem label="创建时间">{{timeUtil.formatDateTime(item['addTime'])}}</FormItem>
          <FormItem label="文档数量">{{item['documentCount']}}</FormItem>
        </Card>
      </ListItem>
    </template>
  </List>
  <Pagination :current="query.page" :page-size="query.pageSize" :total="total" @change="listKnowledgeBase"></Pagination>
  <Drawer title="创建知识库" :open="createKBDrawerOpen" @close="_=>{createKBDrawerOpen=false;}">
    <Form>
      <FormItem label="名称">
        <Input v-model:value="newKnowledgeBase.name"/>
      </FormItem>
      <FormItem label="描述">
        <Textarea v-model:value="newKnowledgeBase.description"/>
      </FormItem>
      <FormItem label="嵌入模型">
        <Select v-model:value="newKnowledgeBase['embeddingModel']" :options="embeddingModels"></Select>
      </FormItem>
    </Form>
    <Button type="primary" @click="createKnowledgeBase">创建</Button>
  </Drawer>
</template>

<style scoped>

</style>