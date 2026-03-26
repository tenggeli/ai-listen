import { createHttpClient } from "@listen/app-shared";

const ADMIN_TOKEN_KEY = "listen:admin:access_token";
const defaultBaseURL = "http://localhost:8080/api/v1/admin";

function getBaseURL(): string {
  return (import.meta.env.VITE_ADMIN_API_BASE_URL as string | undefined) ?? defaultBaseURL;
}

export function getAdminAccessToken(): string {
  return uni.getStorageSync(ADMIN_TOKEN_KEY) || "";
}

export function setAdminAccessToken(token: string) {
  uni.setStorageSync(ADMIN_TOKEN_KEY, token);
}

export function clearAdminAccessToken() {
  uni.removeStorageSync(ADMIN_TOKEN_KEY);
}

export const adminApiClient = createHttpClient({
  baseURL: getBaseURL(),
  getAccessToken: getAdminAccessToken,
  onUnauthorized: clearAdminAccessToken
});
