<script setup>
import {
  Card, Layout, LayoutContent, LayoutSider, Menu, MenuItem
} from "ant-design-vue";
import {onMounted, ref} from "vue";
import knowledgeBaseAPI from '../../api/knowledgeBase.js';
import {useRoute} from "vue-router";
import KnowledgeBaseDocumentList from "./KnowledgeBaseDocumentList.vue";
import KnowledgeBaseSearch from "./KnowledgeBaseSearch.vue";
const route = useRoute()
const id = route.params['id'];
const knowledgeBaseDetail = ref({});
const selectedKeys = ref(['1']);

onMounted(()=>{
  getKnowledgeBaseDetail();
});

function getKnowledgeBaseDetail() {
  knowledgeBaseAPI.detail(id).then((response) => {
    knowledgeBaseDetail.value = response.data;
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
      <Card style="min-height: 100vh;" :title="knowledgeBaseDetail.name">
        <KnowledgeBaseDocumentList v-if="selectedKeys[0]==='1'"/>
        <KnowledgeBaseSearch v-else-if="selectedKeys[0]==='2'"/>
      </Card>
    </LayoutContent>
  </Layout>
</template>

<style scoped>

</style>