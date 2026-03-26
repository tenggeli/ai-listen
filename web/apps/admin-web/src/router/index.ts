import { createRouter, createWebHistory } from "vue-router";
import { useAdminAuthStore } from "../stores/auth";
import DashboardView from "../views/DashboardView.vue";
import LoginView from "../views/LoginView.vue";
import ContentGovernanceView from "../views/ContentGovernanceView.vue";
import ComplaintHandlingView from "../views/ComplaintHandlingView.vue";

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: "/",
      redirect: "/dashboard"
    },
    {
      path: "/login",
      name: "login",
      component: LoginView,
      meta: { guestOnly: true }
    },
    {
      path: "/dashboard",
      name: "dashboard",
      component: DashboardView,
      meta: { requiresAuth: true }
    },
    {
      path: "/governance/content",
      name: "content-governance",
      component: ContentGovernanceView,
      meta: { requiresAuth: true }
    },
    {
      path: "/governance/complaints",
      name: "complaint-handling",
      component: ComplaintHandlingView,
      meta: { requiresAuth: true }
    }
  ]
});

router.beforeEach(async (to) => {
  const authStore = useAdminAuthStore();

  if (to.meta.requiresAuth && !authStore.isLoggedIn) {
    return { path: "/login", query: { redirect: to.fullPath } };
  }

  if (to.meta.guestOnly && authStore.isLoggedIn) {
    return { path: "/dashboard" };
  }

  if (to.meta.requiresAuth && authStore.isLoggedIn && authStore.adminUser == null) {
    try {
      await authStore.fetchMe();
    } catch (error) {
      authStore.clearSession();
      return { path: "/login", query: { redirect: to.fullPath } };
    }
  }

  return true;
});

export default router;
