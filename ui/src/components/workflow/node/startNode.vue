<script setup>
import { Handle, Position } from '@vue-flow/core';
import {Card} from "ant-design-vue";
import NodeVariableDisplay from "./nodeVariableDisplay.vue";

defineProps(['id', 'type', 'data']);
</script>

<template>
	<div>
		<Handle type="target" :position="Position.Right" :connectable="true"></Handle>
    <Card title="开始" :hoverable="true" body-style="padding:10px" head-style="padding:10px;">
      <node-variable-display :has-input="true"
                             :has-output="false"
                             :input-variables="data['startNodeData'].inputVariables"></node-variable-display>
    </Card>
	</div>
</template>

<style scoped></style>