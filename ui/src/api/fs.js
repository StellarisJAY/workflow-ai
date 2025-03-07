import api from "./request"

export default {
    upload: function(form) {
        return api.postForm("/fs/upload", form);
    },
    batchUpload: function(form) {
        return api.postForm("/fs/upload-batch", form);
    },
    fileSrc: function(id) {
        return api.baseURL + "/fs/download/" + id;
    }
}