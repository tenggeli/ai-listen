import { userRoutes } from "./routes";

export function toUserLogin() {
  uni.navigateTo({ url: userRoutes.login });
}

export function toUserProfile() {
  uni.navigateTo({ url: userRoutes.profile });
}
