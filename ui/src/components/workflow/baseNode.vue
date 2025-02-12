<script setup>
import { DeleteOutlined } from '@ant-design/icons-vue';
import { useVueFlow } from '@vue-flow/core';
import { NodeToolbar } from '@vue-flow/node-toolbar';
import { Button } from 'ant-design-vue';

defineProps(['id', 'type', 'data', 'showControls']);
const { removeNodes } = useVueFlow();
function removeNode(id) {
	removeNodes(id);
}
</script>

<template>
	<NodeToolbar :is-visible="showControls">
		<Button type="primary" danger shape="circle" @click="removeNode(id)">
			<DeleteOutlined />
		</Button>
	</NodeToolbar>
</template>

<style scoped></style>