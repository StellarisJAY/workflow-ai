<script setup>
import {
  Form, FormItem, Button,
  Textarea, Card, Row, Col, Slider, InputNumber, message, Divider,
  Collapse, CollapsePanel, Select
} from "ant-design-vue";
import {ref} from "vue";
import knowledgeBaseAPI from "../../api/knowledgeBase.js";
import {useRoute} from "vue-router";
import {kbSearchTypes} from "../../api/const.js";

const route = useRoute()
const kbId = route.params["id"];

const searchRequest = ref({
  scoreThreshold: 0.5,
  count: 10,
  input: "",
  kbId: kbId
});
const searchResult = ref([]);
const activeKey = ref("0");

const searchOption = ref("similarity");

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

function downloadFile(id) {
  knowledgeBaseAPI.downloadFile(id);
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
          <FormItem label="相似度阈值" v-if="searchOption==='similarity'">
            <Slider v-model:value="searchRequest.scoreThreshold" :min="0.01" :max="0.99" :step="0.01"/>
          </FormItem>
          <FormItem label="最大返回数量">
            <InputNumber v-model:value="searchRequest.count" :min="1"/>
          </FormItem>
          <FormItem label="搜索内容">
            <Textarea v-model:value="searchRequest.input" style="height: 300px"/>
          </FormItem>
        </Form>
        <Button v-if="searchOption==='similarity'" type="primary" @click="similaritySearch">搜索</Button>
        <Button v-if="searchOption==='fulltext'" type="primary" @click="fulltextSearch">搜索</Button>
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
          <Textarea v-model:value="item['content']" style="height: 300px"/>
          <Divider/>
        </div>
      </Card>
    </Col>
  </Row>


</template>

<style scoped>

</style>