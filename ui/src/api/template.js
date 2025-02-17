import api from "./request"

export default {
    createTemplate: function(data) {
        return api.post("/template/create", data);
    },
    getTemplate: function(id) {
        return api.get("/template/detail/" + id);
    },
    listTemplate: function(query) {
        return api.get("/template/list", query);
    },
    deleteTemplate: function(id) {
        return api.delete("/template/" + id);
    },
    updateTemplate: function(template) {
        return api.put("/template/update", template);
    },
    getInputVariables: function(templateId) {
        return api.get("/template/start-variables/" + templateId);
    }
}