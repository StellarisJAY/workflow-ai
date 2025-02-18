<script setup>
import {useRoute, useRouter} from "vue-router";
import {ref} from "vue";
import {Layout, LayoutSider, Menu, MenuItem, LayoutContent, Card} from "ant-design-vue";
import TemplateList from "../components/template/templateList.vue";
import InstanceList from "../components/instance/instanceList.vue";
import LlmList from "../components/llm/llmList.vue";

const router = useRouter();
const route = useRoute();
const selectedKeys = ref(['1']);

</script>

<template>
  <Layout style="min-height: 100vh">
    <LayoutSider>
      <Menu v-model:selectedKeys = "selectedKeys" theme="dark" mode="inline">
        <MenuItem key="1">
          <span>流程配置</span>
        </MenuItem>
        <MenuItem key="2">
          <span>模型配置</span>
        </MenuItem>
        <MenuItem key="3">
          <span>知识库</span>
        </MenuItem>
        <MenuItem key="4">
          <span>聊天</span>
        </MenuItem>
        <MenuItem key="5">
          <span>流程实例</span>
        </MenuItem>
      </Menu>
    </LayoutSider>
    <LayoutContent has-sider>
      <Card style="min-height: 100vh">
        <template-list v-if="selectedKeys[0] === '1'"/>
        <llm-list v-if="selectedKeys[0] === '2'"></llm-list>
        <instance-list v-else-if="selectedKeys[0] === '5'"/>
      </Card>
    </LayoutContent>
  </Layout>
</template>

<style scoped>
#components-layout-demo-side .logo {
  height: 32px;
  margin: 16px;
  background: rgba(255, 255, 255, 0.3);
}

.site-layout .site-layout-background {
  background: #fff;
}
[data-theme='dark'] .site-layout .site-layout-background {
  background: #141414;
}
</style>