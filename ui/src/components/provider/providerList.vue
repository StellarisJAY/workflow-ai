<script setup>
import {useRoute, useRouter} from "vue-router";
import {onMounted, ref} from "vue";
import {
  Table,
  Form,
  FormItem,
  Input,
  Select,
  PageHeader,
  Button,
  Drawer,
  CollapsePanel,
  Collapse,
  Card,
  message
} from "ant-design-vue";
import providerAPI from "../../api/provider.js";

const providerModelColumns = [
  {title: "模型", dataIndex: "modelName", key: "modelName" },
  {title: "类型", dataIndex: "modelTypeName", key: "modelTypeName" },
  {title: "操作", dataIndex: "operation", key: "operation" }
]

const modelTypes = [
  {label: "大模型", value: "llm"},
  {label: "文本嵌入", value: "embedding"},
  {label: "图像理解", value: "image_understanding"},
]
const router = useRouter();
const route = useRoute();

const providerList = ref([]);
const query = ref({
  page: 1,
  pageSize: 10,
  paged: false,
});

const newProviderModel = ref({});
const addProviderModelDrawerOpen = ref(false);

function listProviders() {
  providerAPI.listProviders(query.value).then(resp => {
    providerList.value = resp.data;
  });
}

function listProviderModels(providerId) {
  providerAPI.listProviderModels({providerId: providerId, paged: false}).then(resp=>{
    const p = providerList.value.find(p=>p['id'] === providerId);
    if (p) p.models = resp.data;
    p.models.forEach(model=>{
      const type = modelTypes.find(t=>t.value === model.modelType);
      if (type) model.modelTypeName = type.label;
    })
  });
}

function openAddProviderModelDrawer(providerId) {
  newProviderModel.value = {providerId: providerId};
  addProviderModelDrawerOpen.value = true;
}

function addProviderModel() {
  if (!newProviderModel.value['modelType']) {
    message.error("请选择模型类型");
    return;
  }
  if (!newProviderModel.value['modelName']) {
    message.error("请输入模型名称");
    return;
  }
  providerAPI.addProviderModel(newProviderModel.value).then(resp => {
    message.success("添加成功");
    addProviderModelDrawerOpen.value = false;
    listProviderModels(newProviderModel.value.providerId);
  });
}

onMounted(_=>{
  listProviders();
});
</script>

<template>
  <PageHeader title="模型供应商">
    <template #extra>
      <Button>添加供应商</Button>
    </template>
  </PageHeader>
  <Card v-for="provider in providerList" :key="provider.id" :title="provider.name" :hoverable="true">
    <template #extra>
      <Button size="small" @click="openAddProviderModelDrawer(provider.id)">添加模型</Button>
      <Button size="small">设置</Button>
    </template>
    <Collapse @change="_=>listProviderModels(provider.id)">
      <CollapsePanel header="模型列表" >
        <Table :columns="providerModelColumns" :data-source="provider.models" :pagination="false"></Table>
      </CollapsePanel>
    </Collapse>
  </Card>

  <Drawer :open="addProviderModelDrawerOpen"
          @close="_=>{addProviderModelDrawerOpen=false;}"
          :destory-on-close="true" title="添加模型">
    <Form>
      <FormItem label="模型类型">
        <Select v-model:value="newProviderModel.modelType" :options="modelTypes"/>
      </FormItem>
      <FormItem label="模型名称">
        <Input v-model:value="newProviderModel.modelName"/>
      </FormItem>
    </Form>
    <Button type="primary" @click="addProviderModel">确认</Button>
  </Drawer>
</template>

<style scoped>
</style>