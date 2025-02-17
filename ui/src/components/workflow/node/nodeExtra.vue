<script setup>
import { EllipsisOutlined }
  from '@ant-design/icons-vue';
import { useVueFlow } from '@vue-flow/core';
import { Dropdown, Menu, MenuItem } from 'ant-design-vue';
import NodeStatusTag from "./nodeStatusTag.vue";

const props = defineProps(['id', 'type', 'data', 'status', 'editable']);
const { removeNodes, nodesDraggable, nodesConnectable} = useVueFlow();
function removeNode(id) {
  removeNodes(id);
}
</script>

<template>
  <Dropdown v-if="nodesDraggable && nodesConnectable && editable">
    <a class="ant-dropdown-link" @click.prevent>
      <EllipsisOutlined />
    </a>
    <template #overlay>
      <Menu>
        <MenuItem @click="removeNode(id)">
          删除
        </MenuItem>
        <MenuItem>
          复制
        </MenuItem>
      </Menu>
    </template>
  </Dropdown>
  <NodeStatusTag v-if="!nodesDraggable && !nodesConnectable" :status="status"/>
</template>

<style scoped></style>