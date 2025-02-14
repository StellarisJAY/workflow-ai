<script setup>
import editor from './components/workflow/editor.vue';
import templateAPI from './api/template.js';
import {ref} from "vue";

const template = ref({});

templateAPI.getTemplate("1890320070542888960").then(resp=>{
  template.value = resp.data;
});
</script>

<template>
    <editor :is-new-template="false" :template="template"></editor>
</template>

<style scoped>
</style>
