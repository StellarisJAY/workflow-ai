<script setup>
import {useRoute, useRouter} from "vue-router";
import viewer from "../components/workflow/viewer.vue";

const router = useRouter();
const route = useRoute();

const workflowId = route.params['id'];
</script>

<template>
  <viewer :workflow-id="workflowId"></viewer>
</template>

<style scoped>

</style>