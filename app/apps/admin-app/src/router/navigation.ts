import { adminRoutes } from "./routes";

export function toAdminLogin() {
  uni.navigateTo({ url: adminRoutes.login });
}
