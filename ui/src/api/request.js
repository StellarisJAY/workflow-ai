import axios from "axios";

const baseURL = "http://localhost:8080/api/v1";
const v1 = axios.create({
    baseURL: baseURL,
});

v1.interceptors.request.use(config => {
    return config
});

v1.interceptors.response.use(r => {
    const resp = r["data"];
    if (resp && resp.code === 200) {
        return resp;
    }
    return Promise.reject(resp);
});

function errorHandler(err) {
    return Promise.reject(err);
}

const api = {
    axios: v1,
    baseURL: baseURL,
    cli: axios.create({baseURL: baseURL}),
    get(path, queryParams) {
        if (queryParams) {
            let params = [];
            for (let key in queryParams) {
                params.push(key + "=" + queryParams[key]);
            }
            path += "?" + params.join("&");
        }
        return this.axios.get(path).catch(err => errorHandler(err));
    },
    post(path, data) {
        return this.axios.post(path, data).catch(err => errorHandler(err));
    },
    put(path, data) {
        return this.axios.put(path, data).catch(err => errorHandler(err));
    },
    delete(path) {
        return this.axios.delete(path).catch(err => errorHandler(err));
    },
    postForm(path, form) {
        return this.axios.postForm(path, form).catch(err => errorHandler(err));
    }
}

export default api