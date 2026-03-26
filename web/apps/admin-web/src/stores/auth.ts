import { defineStore } from "pinia";
import { computed, ref } from "vue";
import { adminLogin, adminLogout, adminMe, type AdminUser } from "../api/admin";
import { clearAdminAccessToken, getAdminAccessToken, setAdminAccessToken } from "../api/client";

export const useAdminAuthStore = defineStore("admin-auth", () => {
  const accessToken = ref(getAdminAccessToken());
  const adminUser = ref<AdminUser | null>(null);
  const permissions = ref<string[]>([]);
  const loading = ref(false);

  const isLoggedIn = computed(() => accessToken.value.trim() !== "");

  async function login(username: string, password: string) {
    loading.value = true;
    try {
      const data = await adminLogin(username, password);
      accessToken.value = data.accessToken;
      adminUser.value = data.adminUser;
      permissions.value = data.permissions;
      setAdminAccessToken(data.accessToken);
      return data;
    } finally {
      loading.value = false;
    }
  }

  async function fetchMe() {
    if (!isLoggedIn.value) {
      return null;
    }
    const data = await adminMe();
    adminUser.value = data.adminUser;
    permissions.value = data.permissions;
    return data;
  }

  async function logout() {
    if (isLoggedIn.value) {
      try {
        await adminLogout();
      } catch (error) {
        // Ignore network failures for logout and clear local state anyway.
      }
    }
    clearSession();
  }

  function clearSession() {
    accessToken.value = "";
    adminUser.value = null;
    permissions.value = [];
    clearAdminAccessToken();
  }

  function hasPermission(permission: string) {
    return permissions.value.includes(permission);
  }

  return {
    accessToken,
    adminUser,
    permissions,
    loading,
    isLoggedIn,
    login,
    fetchMe,
    logout,
    clearSession,
    hasPermission
  };
});
