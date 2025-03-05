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
    },
    deleteFile: function(id) {
        return api.delete("/knowledgeBase/file/"+id);
    },
    getFileProcessOptions: function(id) {
        return api.get("/knowledgeBase/process/options/" + id);
    },
    updateFileProcessOptions: function(options) {
        return api.put("/knowledgeBase/process/options", options);
    },
    startFileProcessing: function(id) {
        return api.post("/knowledgeBase/process/start/" + id);
    },
    similaritySearch: function(request) {
        return api.post("/knowledgeBase/similarity-search", request);
    },
    fullTextSearch: function(request) {
        return api.post("/knowledgeBase/fulltext-search", request);
    },
    downloadFile: function(id) {
        api.cli.get("/knowledgeBase/download/" + id).then(resp => {
            const _res = resp.data
            let blob = new Blob([_res]);
            let downloadElement = document.createElement("a");
            let href = window.URL.createObjectURL(blob); //创建下载的链接
            downloadElement.href = href;
            downloadElement.download = resp.headers["filename"]; //下载后文件名
            document.body.appendChild(downloadElement);
            downloadElement.click(); //点击下载
            document.body.removeChild(downloadElement); //下载完成移除元素
            window.URL.revokeObjectURL(href); //释放掉blob对象
        });
    },
    listChunks: function(query) {
        return api.get("/knowledgeBase/chunks",query);
    }
}