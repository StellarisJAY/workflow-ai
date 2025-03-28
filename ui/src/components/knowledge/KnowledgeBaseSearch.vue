<script setup>
import {
  Form, FormItem, Button,
  Textarea, Card, Row, Col, Slider, InputNumber, message, Divider,
  Collapse, CollapsePanel, Select, Switch
} from "ant-design-vue";
import {onMounted, ref} from "vue";
import knowledgeBaseAPI from "../../api/knowledgeBase.js";
import {useRoute} from "vue-router";
import {kbSearchTypes} from "../../api/const.js";
import provider from "../../api/provider.js";

const route = useRoute()
const kbId = route.params["id"];

const weightSliderMarks = ref({
  0: "语义 0.5",
  1: "0.5 关键词"
});

const rerankModels = ref([]);
const searchRequest = ref({
  scoreThreshold: 0.5,
  count: 10,
  input: "",
  kbId: kbId,
  hybridSearchOption: {
    weightedRerank: false,
    denseWeight: 0.5,
    sparseWeight: 0.5,
    rerankModelId: ""
  },
});
const searchResult = ref([]);
const activeKey = ref("0");

const searchOption = ref("similarity");

onMounted(()=>{
  listRerankModels();
});

function similaritySearch() {
  knowledgeBaseAPI.similaritySearch(searchRequest.value).then((resp) => {
    searchResult.value = resp.data;
  }).catch((err) => {
    message.error("搜索失败");
  })
}

function fulltextSearch() {
  knowledgeBaseAPI.fullTextSearch(searchRequest.value).then((resp) => {
    searchResult.value = resp.data;
  }).catch((err) => {
    message.error("搜索失败");
  })
}

function hybridSearch() {
  knowledgeBaseAPI.hybridSearch(searchRequest.value).then((resp) => {
    searchResult.value = resp.data;
  }).catch((err) => {
    message.error("搜索失败");
  })
}

function downloadFile(id) {
  knowledgeBaseAPI.downloadFile(id);
}

function onWeightSliderChange(ev) {
  searchRequest.value.hybridSearchOption.sparseWeight = Math.round((1 - ev)*10) / 10;
  searchRequest.value.hybridSearchOption.denseWeight = ev;
  weightSliderMarks.value[0] = searchRequest.value.hybridSearchOption.denseWeight+" 语义";
  weightSliderMarks.value[1] = searchRequest.value.hybridSearchOption.sparseWeight+" 关键词";
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
  <Row>
    <Col :span="8">
      <Card style="height:80vh">
        <Form>
          <FormItem label="搜索方式">
            <Select :options="kbSearchTypes" v-model:value="searchOption"/>
          </FormItem>
          <FormItem label="相似度阈值" v-if="searchOption==='similarity' || searchOption==='hybrid'">
            <Slider v-model:value="searchRequest.scoreThreshold" :min="0.01" :max="0.99" :step="0.01"/>
          </FormItem>
          <FormItem label="最大返回数量">
            <InputNumber v-model:value="searchRequest.count" :min="1"/>
          </FormItem>
          <FormItem label="权重排序" v-if="searchOption==='hybrid'">
            <Switch v-model:checked="searchRequest.hybridSearchOption.weightedRerank"/>
          </FormItem>
          <FormItem label="混合搜索权重" v-if="searchOption==='hybrid' && searchRequest.hybridSearchOption.weightedRerank">
            <Slider :marks="weightSliderMarks" :value="searchRequest.hybridSearchOption.denseWeight" :min="0" :max="1" :step="0.1" @change="onWeightSliderChange">
              <template #mark="{label}">
                {{label}}
              </template>
            </Slider>
          </FormItem>
          <FormItem label="排序模型" v-else-if="searchOption==='hybrid'">
            <Select v-model:value="searchRequest.hybridSearchOption.rerankModelId" :options="rerankModels"/>
          </FormItem>
          <FormItem label="搜索内容">
            <Textarea v-model:value="searchRequest.input" style="height: 300px"/>
          </FormItem>
        </Form>
        <Button v-if="searchOption==='similarity'" type="primary" @click="similaritySearch">搜索</Button>
        <Button v-if="searchOption==='fulltext'" type="primary" @click="fulltextSearch">搜索</Button>
        <Button v-if="searchOption==='hybrid'" type="primary" @click="hybridSearch">搜索</Button>
      </Card>
    </Col>
    <Col :span="16">
      <Card style="max-height:80vh;height:80vh;overflow: auto">
        <Collapse v-model:activeKey="activeKey">
          <CollapsePanel key="1" header="选中的文档">
            <div v-for="file in searchResult['files']">
              <a @click="downloadFile(file['id'])">{{file['name']}}</a>
            </div>
          </CollapsePanel>
        </Collapse>
        <div v-for="item in searchResult['documents']">
          <p style="color:gray">相似度得分:{{item['score']}}</p>
          <div style="height: 300px; width: 100%; overflow:auto">
            <p>{{item['content']}}</p>
          </div>
          <Divider/>
        </div>
      </Card>
    </Col>
  </Row>


</template>

<style scoped>

</style>