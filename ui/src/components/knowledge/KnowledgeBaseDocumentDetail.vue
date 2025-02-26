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
const searchResult = ref([]);
</script>

<template>
  <Card style="max-height:80vh;height:80vh;overflow: auto">
    <div v-for="item in chunks">
      <Textarea v-model:value="item['content']" style="height: 300px"/>
      <Divider/>
    </div>
  </Card>
</template>

<style scoped>

</style>