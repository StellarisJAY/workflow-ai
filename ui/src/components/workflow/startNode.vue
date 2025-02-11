<script setup>
import { Handle, Position } from '@vue-flow/core';

defineProps(['id', 'type', 'data']);
</script>

<template>
	<div>
		<span>开始</span>
		<Handle type="target" :position="Position.Right" :connectable="true"></Handle>
	</div>
</template>

<style scoped></style>