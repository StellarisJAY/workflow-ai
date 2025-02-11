<script setup>
import { EdgeLabelRenderer, getBezierPath, SimpleBezierEdge, useVueFlow } from '@vue-flow/core';
import { computed } from 'vue';

const props = defineProps(['id', 'sourceX', 'sourceY', 'targetX', 'targetY', 'sourcePosition', 'targetPosition', 'markerStart']);
const {removeEdges} = useVueFlow();
const path = computed(_=>getBezierPath(props));
</script>

<template>
	<SimpleBezierEdge :source-x="sourceX" :source-y="sourceY" :target-x="targetX" :target-y="targetY"
		:source-position="sourcePosition" :target-position="targetPosition" :marker-start="markerStart">
	</SimpleBezierEdge>
	<EdgeLabelRenderer>
		<div class="nodrag nopan" :style="{
        pointerEvents: 'all',
        position: 'absolute',
        transform: `translate(-50%, -50%) translate(${path[1]}px,${path[2]}px)`,
      }">
			<button class="edgebutton" @click="removeEdges(id)">×</button>
		</div>
	</EdgeLabelRenderer>
</template>


<style scoped>
.edgebutton {
  border-radius: 999px;
  cursor: pointer;
}

.edgebutton:hover {
  box-shadow: 0 0 0 2px pink, 0 0 0 4px #f05f75;
}
</style>