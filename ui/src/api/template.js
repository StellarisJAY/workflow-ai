import api from "./request"

export default {
    createTemplate: function(data) {
        return api.post("/template/create", data);
    },
    getTemplate: function(id) {
        return api.get("/template/detail/" + id);
    }
}