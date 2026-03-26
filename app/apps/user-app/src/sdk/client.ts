import { createHttpClient } from "@listen/app-shared";

const USER_TOKEN_KEY = "listen:user:access_token";
const defaultBaseURL = "http://localhost:8080/api/v1";

function getBaseURL(): string {
  return (import.meta.env.VITE_API_BASE_URL as string | undefined) ?? defaultBaseURL;
}

export function getUserAccessToken(): string {
  return uni.getStorageSync(USER_TOKEN_KEY) || "";
}

export function setUserAccessToken(token: string) {
  uni.setStorageSync(USER_TOKEN_KEY, token);
}

export function clearUserAccessToken() {
  uni.removeStorageSync(USER_TOKEN_KEY);
}

export const userApiClient = createHttpClient({
  baseURL: getBaseURL(),
  getAccessToken: getUserAccessToken,
  onUnauthorized: clearUserAccessToken
});
