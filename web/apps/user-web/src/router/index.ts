import { createRouter, createWebHistory } from "vue-router";
import HomeView from "../views/HomeView.vue";
import LoginView from "../views/LoginView.vue";
import OrdersView from "../views/OrdersView.vue";
import ProfileView from "../views/ProfileView.vue";

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: "/",
      name: "home",
      component: HomeView
    },
    {
      path: "/login",
      name: "login",
      component: LoginView
    },
    {
      path: "/orders",
      name: "orders",
      component: OrdersView
    },
    {
      path: "/profile",
      name: "profile",
      component: ProfileView
    }
  ]
});

export default router;
