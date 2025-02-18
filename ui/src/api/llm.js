import api from "./request"

export default {
    listModels: function(query) {
        return api.get("/model/list", query);
    },
    getModelDetail: function(id) {
        return api.get("/model/detail/"+id);
    },
    createModel: function(data) {
        return api.post("/model/create",data);
    }
}