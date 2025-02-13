<script setup>
import { Handle, Position } from '@vue-flow/core';
import {Card} from "ant-design-vue";

defineProps(['id', 'type', 'data']);
</script>

<template>
	<div>
		<Handle type="source" :position="Position.Left" :connectable="true"></Handle>
    <Card title="结束" :hoverable="true">
    </Card>
	</div>
</template>

<style scoped></style>