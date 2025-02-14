import api from "./request"

export default {
    start: function(data) {
        return api.post("/workflow/start", data);
    },
    getOutputs: function(workflowId) {
        return api.get("/workflow/outputs/"+workflowId);
    }
}