<script setup>
import { Handle, Position } from '@vue-flow/core';

defineProps(['id', 'type', 'data']);
</script>

<template>
	<div>
		<Handle type="source" :position="Position.Right"></Handle>
		<div>{{ data.title }}</div>
		<Handle type="target" :position="Position.Left"></Handle>
	</div>
</template>

<style scoped></style>