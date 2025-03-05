<script setup>
import {
  Textarea, Divider, Pagination
} from "ant-design-vue";
import {onMounted, ref} from "vue";
import knowledgeBaseAPI from "../../api/knowledgeBase.js";
import {useRoute} from "vue-router";

const props = defineProps(["fileId"]);

const route = useRoute()
const kbId = route.params["id"];
const chunks = ref([]);

const query = ref({
  page: 1,
  pageSize: 10,
  kbId: kbId,
  fileId: props.fileId,
});
const total = ref(0);

onMounted(()=>{
  listChunks();
});

function listChunks() {
  knowledgeBaseAPI.listChunks(query.value).then((resp)=>{
    chunks.value = resp.data;
    total.value = resp.total;
  });
}
</script>

<template>
  <div style="height: 90%; overflow: auto;">
    <div v-for="item in chunks">
      <Textarea v-model:value="item['content']" style="height: 300px"/>
      <Divider/>
    </div>
  </div>

  <Pagination v-model:current="query.page"
              :page-size="query.pageSize"
              :total="total"
              @change="listChunks"
              :show-size-changer="false"/>
</template>

<style scoped>

</style>