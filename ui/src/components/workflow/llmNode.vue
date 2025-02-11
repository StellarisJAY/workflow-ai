<script setup>
import { Handle, Position } from '@vue-flow/core';

defineProps(['id', 'type', 'data']);
</script>

<template>
	<div>
		<p>大模型</p>
		<Handle type="source" :position="Position.Left"></Handle>
		<Handle type="target" :position="Position.Right"></Handle>
	</div>
</template>

<style scoped></style>