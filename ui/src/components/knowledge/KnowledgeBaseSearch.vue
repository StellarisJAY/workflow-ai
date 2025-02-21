<script setup>
import {
  Form, FormItem, Button,
  Textarea, Card, Row, Col, Slider, InputNumber, message, Divider,
  Collapse, CollapsePanel, List, ListItem
} from "ant-design-vue";
import {ref} from "vue";
import knowledgeBaseAPI from "../../api/knowledgeBase.js";
import {useRoute} from "vue-router";
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
function similaritySearch() {
  knowledgeBaseAPI.similaritySearch(searchRequest.value).then((resp) => {
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
          <FormItem label="相似度阈值">
            <Slider v-model:value="searchRequest.scoreThreshold" :min="0" :max="1" :step="0.1"/>
          </FormItem>
          <FormItem label="最大返回数量">
            <InputNumber v-model:value="searchRequest.count" :min="1"/>
          </FormItem>
          <FormItem label="搜索内容">
            <Textarea v-model:value="searchRequest.input" style="height: 300px"/>
          </FormItem>
        </Form>
        <Button type="primary" @click="similaritySearch">检索</Button>
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