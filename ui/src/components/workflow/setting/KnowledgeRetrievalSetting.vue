<script setup>
import VariableTable from "./VariableTable.vue";
import CommonSetting from "./CommonSetting.vue";
import {Form, FormItem, InputNumber, Slider, Select, Switch} from "ant-design-vue";
import {kbSearchTypes} from "../../../api/const.js";
import {onBeforeMount, onMounted, ref} from "vue";
import knowledgeBaseAPI from "../../../api/knowledgeBase.js";
import provider from "../../../api/provider.js";

const props = defineProps(['node']);

const kbList = ref([]);
const kbOptions = ref([]);
const rerankModels = ref([]);
const weightSliderMarks = ref({
  0: "语义 0.5",
  1: "0.5 关键词"
});

onBeforeMount(()=>{
  if (!props.node.data['retrieveKnowledgeBaseNodeData'].hybridSearchOption) {
    props.node.data['retrieveKnowledgeBaseNodeData'].hybridSearchOption = {
      weightedRerank: false,
      sparseWeight: 0.5,
      denseWeight: 0.5,
      rerankModelId: ""
    };
  }
})

onMounted(()=>{

  knowledgeBaseAPI.list({paged: false}).then(resp=>{
    kbList.value = resp.data;
    const options = [];
    kbList.value.forEach(item=>{
      options.push({label: item.name, value: item.id});
    });
    kbOptions.value = options;
  });
  listRerankModels();
});

function onWeightSliderChange(ev) {
  const nodeData = props.node.data["retrieveKnowledgeBaseNodeData"];
  nodeData.hybridSearchOption.sparseWeight = Math.round((1 - ev)*10) / 10;
  nodeData.hybridSearchOption.denseWeight = ev;
  weightSliderMarks.value[0] = nodeData.hybridSearchOption.denseWeight+" 语义";
  weightSliderMarks.value[1] = nodeData.hybridSearchOption.sparseWeight+" 关键词";
}


function listRerankModels() {
  provider.listProviderModels({modelType: "text_rerank"}).then(resp=>{
    const models = resp.data;
    const options = [];
    models.forEach(model=>{
      options.push({label: model['providerName'] + "/" + model.modelName, value: model.id});
    });
    rerankModels.value = options;
  });
}

</script>

<template>
  <CommonSetting :node="node"/>
  <Form>
    <FormItem label="知识库">
      <Select :options="kbOptions" v-model:value="node.data['retrieveKnowledgeBaseNodeData']['kbId']"/>
    </FormItem>
    <FormItem label="搜索类型">
      <Select :options="kbSearchTypes" v-model:value="node.data['retrieveKnowledgeBaseNodeData']['searchType']"/>
    </FormItem>
    <FormItem v-if="node.data['retrieveKnowledgeBaseNodeData']['searchType'] !== 'fulltext'" label="相似度阈值">
      <Slider :min="0.01" :max="0.99" :step="0.01" v-model:value="node.data['retrieveKnowledgeBaseNodeData']['similarityThreshold']"/>
    </FormItem>
    <FormItem label="最大返回文档数量">
      <InputNumber :min="1" v-model:value="node.data['retrieveKnowledgeBaseNodeData']['count']"/>
    </FormItem>
    <FormItem label="使用权重排序" >
      <Switch v-model:checked="node.data['retrieveKnowledgeBaseNodeData'].hybridSearchOption.weightedRerank"/>
    </FormItem>
    <div v-if="node.data['retrieveKnowledgeBaseNodeData']['searchType']==='hybrid'">
      <FormItem label="混合搜索权重" v-if="node.data['retrieveKnowledgeBaseNodeData'].hybridSearchOption.weightedRerank">
        <Slider :marks="weightSliderMarks"
                :value="node.data['retrieveKnowledgeBaseNodeData'].denseWeight"
                :min="0"
                :max="1"
                :step="0.1"
                @change="onWeightSliderChange">
          <template #mark="{label, point}">
            {{label}}
          </template>
        </Slider>
      </FormItem>
      <FormItem label="排序模型" v-else>
        <Select v-model:value="node.data['retrieveKnowledgeBaseNodeData'].hybridSearchOption.rerankModelId" :options="rerankModels"/>
      </FormItem>
    </div>
  </Form>
  <VariableTable :node-id="node.id"
                 :input-variables="node.data['input']"
                 :output-variables="node.data['output']"
                 :node-data="node.data"
                 :node="node"/>
</template>

<style scoped></style>