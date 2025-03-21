<script setup>
import VariableTable from "./VariableTable.vue";
import CommonSetting from "./CommonSetting.vue";
import {Form, FormItem, InputNumber, Slider, Select} from "ant-design-vue";
import {kbSearchTypes} from "../../../api/const.js";
import {onMounted, ref} from "vue";
import knowledgeBaseAPI from "../../../api/knowledgeBase.js";

const props = defineProps(['node']);

const kbList = ref([]);
const kbOptions = ref([]);

const weightSliderMarks = ref({
  0: "语义 0.5",
  1: "0.5 关键词"
});

onMounted(()=>{
  knowledgeBaseAPI.list({paged: false}).then(resp=>{
    kbList.value = resp.data;
    const options = [];
    kbList.value.forEach(item=>{
      options.push({label: item.name, value: item.id});
    });
    kbOptions.value = options;
  });
});

function onWeightSliderChange(ev) {
  const nodeData = props.node.data["retrieveKnowledgeBaseNodeData"];
  nodeData.sparseWeight = Math.round((1 - ev)*10) / 10;
  nodeData.denseWeight = ev;
  weightSliderMarks.value[0] = nodeData.denseWeight+" 语义";
  weightSliderMarks.value[1] = nodeData.sparseWeight+" 关键词";
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
    <FormItem label="混合搜索权重" v-if="node.data['retrieveKnowledgeBaseNodeData']['searchType']==='hybrid'">
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
  </Form>
  <VariableTable :node-id="node.id"
                 :input-variables="node.data['input']"
                 :output-variables="node.data['output']"
                 :node-data="node.data"
                 :node="node"/>
</template>

<style scoped></style>