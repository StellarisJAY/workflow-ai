<script setup>
import { Handle, Position } from '@vue-flow/core';
import {Card} from "ant-design-vue";
import NodeVariableDisplay from "./nodeVariableDisplay.vue";

defineProps(['id', 'type', 'data']);
</script>

<template>
	<div>
		<Handle type="target" :position="Position.Right" :connectable="true"></Handle>
    <Card title="开始" :hoverable="true">
      <node-variable-display :has-input="true"
                             :has-output="false"
                             :input-variables="data.inputVariables"></node-variable-display>
    </Card>
	</div>
</template>

<style scoped></style>