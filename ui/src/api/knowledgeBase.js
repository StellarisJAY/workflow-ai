import api from "./request"

export default {
    create: function(data) {
        return api.post("/knowledgeBase/create", data);
    },
    list: function(query) {
        return api.get("/knowledgeBase/list", query);
    },
    detail: function(id) {
        return api.get("/knowledgeBase/detail/" + id);
    },
    listDocuments: function(kbId, query) {
        return api.get("/knowledgeBase/files/" + kbId, query);
    },
    upload: function(form) {
        return api.postForm("/knowledgeBase/upload", form);
    }
}