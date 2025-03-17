<script setup>

import {Button, Form, FormItem, Input, message, Spin, Upload} from "ant-design-vue";
import {onMounted, ref} from "vue";
import templateAPI from "../../api/template.js";
import fsAPI from "../../api/fs.js";
import workflowAPI from "../../api/workflow.js";
import {useRouter} from "vue-router";

const props = defineProps(["startingTemplateId"]);

const executeInputOpen = ref(false);
const executeInputVars = ref([]);
const executePending = ref(false);

const router = useRouter();

const fileLists = ref({});

onMounted(() => {
  listInputVariables(props.startingTemplateId);
});


function listInputVariables(id) {
  executePending.value = true;
  templateAPI.getInputVariables(id).then(resp=>{
    executeInputVars.value = resp.data;
    executePending.value = false;
  });
}

function startWorkflow() {
  executePending.value = true;
  const form = new FormData();
  for (let k in fileLists.value) {
    form.append(k, fileLists.value[k][0].originFileObj);
  }
  // 上传文件，获取文件变量的id
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
      templateId: props.startingTemplateId,
    };
    workflowAPI.start(request).then(resp=>{
      message.success("开始执行成功");
      const workflowId = resp.data['workflowId'];
      router.push("/view/"+workflowId);
    }).catch(()=>{
      executePending.value = false;
    });
  }).catch(err=>{
    console.log(err);
    executePending.value = false;
  })

}

function onUploadFileChange(ev, variable) {
  fileLists.value[variable.name] = ev.fileList;
}
</script>

<template>
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
</template>