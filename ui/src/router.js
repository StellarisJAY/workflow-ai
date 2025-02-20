import {createRouter, createWebHistory} from "vue-router";
import Editor from "./pages/editor.vue";
import Home from "./pages/home.vue";
import Viewer from "./pages/viewer.vue";
import KnowledgeBaseDetail from "./components/knowledge/KnowledgeBaseDetail.vue";

const routes = [
    {
        path: "/editor/:id",
        name: "Editor",
        component: Editor,
    },
    {
        path: "/",
        name: "Home",
        component: Home,
    },
    {
        path: "/view/:id",
        name: "View",
        component: Viewer,
    },
    {
        path: "/knowledgeBase/:id",
        name: "KnowledgeBase",
        component: KnowledgeBaseDetail,
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes: routes
});

export default router;