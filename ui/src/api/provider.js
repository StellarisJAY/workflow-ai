import api from "./request"

export default {
    listProviders: function(query) {
        return api.get("/provider/list", query);
    },
    listProviderModels: function(query) {
        return api.get("/provider/model/list", query);
    },
    addProviderModel: function(pm) {
        return api.post("/provider/model/create", pm);
    }
}