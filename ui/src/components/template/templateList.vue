<script setup>
import {useRoute, useRouter} from "vue-router";
import {onMounted, ref} from "vue";
import templateAPI from "../../api/template.js";
import {Card, List, ListItem, Pagination, Dropdown, Menu, MenuItem, Popconfirm, message,
  PageHeader, Button, Drawer, Form, FormItem, Input, Upload, Spin} from "ant-design-vue";
import {EditOutlined, PlayCircleOutlined, EllipsisOutlined} from "@ant-design/icons-vue";
import workflowAPI from "../../api/workflow.js";
import fsAPI from "../../api/fs.js";

const router = useRouter();
const route = useRoute();
const templateList = ref([]);
const query = ref({
  page: 1,
  pageSize: 10,
  paged: true,
});
const total = ref(0);
const executeInputOpen = ref(false);
const executeInputVars = ref([]);
const startingTemplateId = ref("");
const executePending = ref(false);

const fileLists = ref({});

function listTemplates() {
  templateAPI.listTemplate(query.value).then(resp=>{
    templateList.value = resp.data;
    total.value = resp.total;
  });
}

function openEditor(id) {
  router.push('/editor/' + id);
}

function deleteTemplate(id) {
  templateAPI.deleteTemplate(id).then(_=>{
    message.success("删除成功")
    listTemplates();
  }).catch(_=>{
    message.error("删除失败");
  });
}

function openNewTemplate() {
  router.push('/editor/new');
}

function openExecuteInputVariableDrawer(id) {
  templateAPI.getInputVariables(id).then(resp=>{
    executeInputVars.value = resp.data;
    startingTemplateId.value = id;
    executeInputOpen.value = true;
  });
}

function startWorkflow() {
  executePending.value = true;
  const form = new FormData();
  for (let k in fileLists.value) {
    form.append(k, fileLists.value[k][0].originFileObj);
  }
  fsAPI.batchUpload(form).then(resp=>{
    const fileIds = resp.data['fileIds'];
    let i = 0;
    for (let k in fileLists.value) {
      const variable = executeInputVars.value.find(item=>item.name === k);
      if (variable) {
        variable.value = fileIds[i];
        i++;
      }
    }
    const input = {};
    executeInputVars.value.forEach((item) => {
      input[item.name] = item.value;
    });
    const request = {
      inputs: input,
      templateId: startingTemplateId.value,
    };
    workflowAPI.start(request).then(resp=>{
      message.success("开始执行成功");
      const workflowId = resp.data['workflowId'];
      router.push("/view/"+workflowId);
    }).catch(()=>{
      executePending.value = false;
    });
  }).catch(err=>{
    executePending.value = false;
  })

}

function onUploadFileChange(ev, variable) {
  fileLists.value[variable.name] = ev.fileList;
}

onMounted(_=>{
  listTemplates();
});
</script>

<template>
  <PageHeader title="应用列表">
    <template #extra>
      <Button type="primary" @click="openNewTemplate">新建</Button>
    </template>
  </PageHeader>
  <List :grid="{ gutter: 4, column: 4 }" :data-source="templateList">
    <template #renderItem="{item}">
      <ListItem>
        <Card :title="item.name">
          <template #actions>
            <PlayCircleOutlined @click="openExecuteInputVariableDrawer(item.id)" />
            <EditOutlined @click="openEditor(item.id)"/>
          </template>
          <template #extra>
            <Dropdown>
              <a class="ant-dropdown-link" @click.prevent>
                <EllipsisOutlined />
              </a>
              <template #overlay>
                <Menu>
                  <MenuItem>
                    <Popconfirm title="确认删除流程模板？"
                                @confirm="deleteTemplate(item.id)"
                                confirm-text="确认"
                                cancel-text="取消">
                      删除
                    </Popconfirm>
                  </MenuItem>
                </Menu>
              </template>
            </Dropdown>
          </template>
          {{item.addTime}}
        </Card>
      </ListItem>
    </template>
  </List>
  <Pagination :current="query.page" :page-size="query.pageSize" :total="total"></Pagination>
  <Drawer title="输入变量" :open="executeInputOpen" @close="_=>{executeInputOpen = false;}">
    <Spin :spinning="executePending">
      <Form>
        <FormItem v-for="variable in executeInputVars" :label="variable.name">
          <Input v-model:value="variable.value" v-if="variable.type === 'string'" placeholder="变量值"/>
          <Upload v-else :file-list="fileLists[variable.name]" :multiple="false"
                  :before-upload="()=>false" @change="ev=>onUploadFileChange(ev, variable)">
            <Button size="small">上传</Button>
          </Upload>
        </FormItem>
      </Form>
      <Button type="primary" @click="startWorkflow">执行流程</Button>
    </Spin>
  </Drawer>
</template>

<style scoped>
</style>