<script setup>
import { Handle, Position } from '@vue-flow/core';
import {FormItem} from "ant-design-vue";

defineProps(['id', 'type', 'data']);
</script>

<template>
	<div>
		<span>开始</span>
		<Handle type="target" :position="Position.Right" :connectable="true"></Handle>
    <FormItem label="输入变量">{{data['inputVariables']}}</FormItem>
	</div>
</template>

<style scoped></style>