<script setup>
import { Handle, Position } from '@vue-flow/core';

defineProps(['id', 'type', 'data']);
</script>

<template>
	<div>
		<p>条件</p>
		<Handle type="source" :position="Position.Left">+</Handle>
		<Handle id="1" type="target" :position="Position.Right" style="top: 10px; bottom: auto;">1</Handle>
	</div>
</template>

<style scoped></style>