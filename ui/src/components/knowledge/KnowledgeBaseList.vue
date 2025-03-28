<script setup>
import {
  Card,
  List,
  ListItem,
  Button,
  Drawer,
  Form,
  PageHeader,
  Pagination, FormItem, Input, Textarea, Select, message, Spin, Menu, MenuItem, Popconfirm, Dropdown
} from "ant-design-vue";
import {onMounted, ref} from "vue";
import knowledgeBaseAPI from '../../api/knowledgeBase.js';
import {EllipsisOutlined, SettingOutlined, UploadOutlined} from "@ant-design/icons-vue";
import timeUtil from "../../util/timeUtil.js";
import {useRouter} from "vue-router";
import provider from "../../api/provider.js";
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
const creating = ref(false);

onMounted(_=>{
  listKnowledgeBase();
});

function openCreateKBDrawer() {
  listEmbeddingModels();
  newKnowledgeBase.value = {name: "", description: ""};
  createKBDrawerOpen.value = true;
}

function listEmbeddingModels() {
  provider.listProviderModels({modelType: "embedding"}).then(resp=>{
    const models = resp.data;
    const options = [];
    models.forEach(model=>{
      options.push({label: model['providerName'] + "/" + model.modelName, value: model.id});
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
  creating.value = true;
  knowledgeBaseAPI.create(newKnowledgeBase.value).then(resp=>{
    message.success("创建成功");
    createKBDrawerOpen.value = false;
    listKnowledgeBase();
    creating.value = false;
  }).catch(err=>{
    console.log(err);
    creating.value = false;
  })
}

function deleteKnowledgeBase(id) {
  knowledgeBaseAPI.delete(id).then(resp=>{
    message.success("删除成功");
    listKnowledgeBase();
  }).catch(err=>{
    message.error(err.message);
  })
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
          <template #extra>
            <Dropdown>
              <a class="ant-dropdown-link" @click.prevent>
                <EllipsisOutlined />
              </a>
              <template #overlay>
                <Menu>
                  <MenuItem>
                    <Popconfirm title="确认删除该知识库？知识库中所有文件将被删除无法恢复"
                                @confirm="deleteKnowledgeBase(item.id)"
                                confirm-text="确认"
                                cancel-text="取消">
                      删除
                    </Popconfirm>
                  </MenuItem>
                </Menu>
              </template>
            </Dropdown>
          </template>
          <template #actions>
            <UploadOutlined />
            <SettingOutlined @click="router.push('/knowledgeBase/'+item.id)"/>
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