import api from "./request"

export default {
    start: function(data) {
        return api.post("/workflow/start", data);
    },
    getOutputs: function(workflowId) {
        return api.get("/workflow/outputs/"+workflowId);
    },
    list: function() {
        return api.get("/workflow/list");
    },
    detail: function(id) {
        return api.get("/workflow/detail/"+id);
    },
    getNodeInstanceDetail: function(workflowId, nodeId) {
        return api.get("/workflow/node/detail", {workflowId: workflowId, nodeId: nodeId});
    }
}