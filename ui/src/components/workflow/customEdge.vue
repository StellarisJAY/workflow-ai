<script setup>
import {BaseEdge, EdgeLabelRenderer, getBezierPath, SimpleBezierEdge, useVueFlow} from '@vue-flow/core';
import {computed, watch} from 'vue';
import {Button} from "ant-design-vue";

const props = defineProps(['id', 'sourceX', 'sourceY', 'targetX', 'targetY',
  'sourcePosition', 'targetPosition', 'markerStart',
  'sourceNode', 'targetNode', 'source', 'target', 'type', 'updatable', 'selected',
  'animated', 'label', 'labelStyle', 'labelShowBg', 'labelBgStyle', 'labelBgPadding',
  'labelBgBorderRadius', 'data', 'events', 'style', 'markerEnd', 'sourceHandleId',
  'targetHandleId', 'interactionWidth']);
const {removeEdges, edgesUpdatable, findEdge} = useVueFlow();
const path = computed(_=>getBezierPath(props));
const edge = findEdge(props.id);
if (edge['passed']) {
  edge.animated = true;
}
</script>

<template>
	<SimpleBezierEdge :source-x="sourceX" :source-y="sourceY" :target-x="targetX" :target-y="targetY"
		:source-position="sourcePosition" :target-position="targetPosition" :marker-end="markerEnd">
	</SimpleBezierEdge>
	<EdgeLabelRenderer>
		<div :style="{
        pointerEvents: 'all',
        position: 'absolute',
        transform: `translate(-50%, -50%) translate(${path[1]}px,${path[2]}px)`}">
			<Button shape="circle" size="small" @click="removeEdges(id)" v-if="edgesUpdatable">x</Button>
		</div>
	</EdgeLabelRenderer>
</template>


<style scoped>
</style>