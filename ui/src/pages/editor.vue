<script setup>
import {useRoute, useRouter} from "vue-router";
import editor from '../components/workflow/editor.vue';
import templateAPI from '../api/template.js';
import {ref} from "vue";

const router = useRouter();
const route = useRoute();

const templateId = route.params['id'];

const isNewTemplate = ref(false);
const template = ref({});
if (templateId === 'new') {
  isNewTemplate.value = true;
}else {
  templateAPI.getTemplate(templateId).then(resp=>{
    template.value = resp.data;
  });
}
</script>

<template>
  <editor :is-new-template="isNewTemplate" :template="template"/>
</template>

<style scoped>

</style>