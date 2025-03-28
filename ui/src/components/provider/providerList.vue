<script setup>
import {useRoute, useRouter} from "vue-router";
import {onMounted, ref} from "vue";
import {
  Button,
  Card,
  Collapse,
  CollapsePanel,
  Drawer,
  Form,
  FormItem,
  Input,
  message,
  PageHeader,
  Select,
  Spin,
  Table
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
  {label: "文本排序", value: "text_rerank"},
]

const providerCodes = [
  {label: "OpenAI", value: "openai"},
  {label: "通义千问", value: "tongyi"},
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

const newProvider = ref({});
const addProviderDrawerOpen = ref(false);
const newProviderLoadingSchema = ref(false);
const providerSchemas = ref(null);
const newProviderCredentials = ref({});

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

function openAddProviderDrawer() {
  newProvider.value = {};
  addProviderDrawerOpen.value = true;
}

function onProviderCodeChange() {
  if (!providerSchemas.value) {
    newProviderLoadingSchema.value = true;
    providerAPI.listProviderSchemas().then(resp => {
      providerSchemas.value = resp.data;
      newProviderLoadingSchema.value = false;
      const schema = providerSchemas.value.find(p=>p.code===newProvider.value.code);
      if (schema) {
        newProviderCredentials.value = schema["credentialSchema"];
      }
    });
  }else {
    const schema = providerSchemas.value.find(p=>p.code===newProvider.value.code);
    if (schema) {
      newProviderCredentials.value = schema["credentialSchema"];
    }
  }
}

function addProvider() {
  for (let k in newProviderCredentials.value) {
    if (!newProviderCredentials.value[k]) {
      message.warn("请填写\""+k);
      return;
    }
  }
  if (!newProvider.value["code"]) {
    message.warn("请选择供应商API类型");
    return;
  }
  if (!newProvider.value["name"]) {
    message.warn("请填写供应商名称");
    return;
  }
  newProvider.value["credentials"] = JSON.stringify(newProviderCredentials.value);
  providerAPI.addProvider(newProvider.value).then(resp => {
    message.success("添加成功");
    addProviderDrawerOpen.value = false;
    listProviders();
  });
}

onMounted(_=>{
  listProviders();
});
</script>

<template>
  <PageHeader title="模型供应商">
    <template #extra>
      <Button @click="openAddProviderDrawer">添加供应商</Button>
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

  <Drawer :open="addProviderDrawerOpen" @close="_=>{addProviderDrawerOpen=false;}" :destroy-on-close="true" title="添加供应商">
    <Form>
      <FormItem label="供应商名称">
        <Input v-model:value="newProvider.name"/>
      </FormItem>
      <FormItem label="API类型">
        <Select v-model:value="newProvider.code" :options="providerCodes" @change="onProviderCodeChange"/>
      </FormItem>
      <Spin :spinning="newProviderLoadingSchema">
        <FormItem v-for="(v, k) in newProviderCredentials" :label="k">
          <Input v-model:value="newProviderCredentials[k]" />
        </FormItem>
      </Spin>
    </Form>
    <Button type="primary" @click="addProvider">确认</Button>
  </Drawer>
</template>

<style scoped>
</style>