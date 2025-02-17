<script setup>
import { EdgeLabelRenderer, getBezierPath, SimpleBezierEdge, useVueFlow } from '@vue-flow/core';
import { computed } from 'vue';
import {Button} from "ant-design-vue";

const props = defineProps(['id', 'sourceX', 'sourceY', 'targetX', 'targetY', 'sourcePosition', 'targetPosition', 'markerStart']);
const {removeEdges, edgesUpdatable} = useVueFlow();
const path = computed(_=>getBezierPath(props));
</script>

<template>
	<SimpleBezierEdge :source-x="sourceX" :source-y="sourceY" :target-x="targetX" :target-y="targetY"
		:source-position="sourcePosition" :target-position="targetPosition" :marker-start="markerStart">
	</SimpleBezierEdge>
	<EdgeLabelRenderer>
		<div :style="{
        pointerEvents: 'all',
        position: 'absolute',
        transform: `translate(-50%, -50%) translate(${path[1]}px,${path[2]}px)`,
      }">
			<Button shape="circle" size="small" @click="removeEdges(id)" v-if="edgesUpdatable">x</Button>
		</div>
	</EdgeLabelRenderer>
</template>


<style scoped>
</style>