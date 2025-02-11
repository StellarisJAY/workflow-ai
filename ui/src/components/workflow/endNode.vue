<script setup>
import { Handle, Position } from '@vue-flow/core';

defineProps(['id', 'type', 'data']);
</script>

<template>
	<div>
		<span>结束</span>
		<Handle type="source" :position="Position.Left" :connectable="true"></Handle>
	</div>
</template>

<style scoped></style>