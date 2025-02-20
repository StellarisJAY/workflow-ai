<script setup>
import {
  Card, Layout, LayoutContent, LayoutSider, Menu, MenuItem,
} from "ant-design-vue";
import {onMounted, ref} from "vue";
import llmAPI from '../../api/llm.js';
import knowledgeBaseAPI from '../../api/knowledgeBase.js';
import {useRoute} from "vue-router";
import KnowledgeBaseDocumentList from "./KnowledgeBaseDocumentList.vue";

const route = useRoute()
const id = route.params['id'];
const embeddingModels = ref([]);
const knowledgeBaseDetail = ref({});
const selectedKeys = ref(['1']);

onMounted(()=>{
  listEmbeddingModels();
  getKnowledgeBaseDetail();
});

function getKnowledgeBaseDetail() {
  knowledgeBaseAPI.detail(id).then((response) => {
    knowledgeBaseDetail.value = response.data;
  });
}

function listEmbeddingModels() {
  llmAPI.listModels({modelType: "embedding"}).then(response => {
    const options = [];
    response.data.forEach((item) => {
      options.push({label: item.name, value: item.id});
    });
    embeddingModels.value = options;
  });
}
</script>

<template>
  <Layout style="min-height: 100vh">
    <LayoutSider  theme="light">
      <Menu v-model:selectedKeys = "selectedKeys" theme="light" mode="inline">
        <MenuItem key="1">
          <span>文档</span>
        </MenuItem>
        <MenuItem key="2">
          <span>检索</span>
        </MenuItem>
        <MenuItem key="3">
          <span>设置</span>
        </MenuItem>
      </Menu>
    </LayoutSider>
    <LayoutContent has-sider>
      <Card style="min-height: 100vh" :title="knowledgeBaseDetail.name">
        <KnowledgeBaseDocumentList v-if="selectedKeys[0]==='1'"/>
      </Card>
    </LayoutContent>
  </Layout>
</template>

<style scoped>

</style>