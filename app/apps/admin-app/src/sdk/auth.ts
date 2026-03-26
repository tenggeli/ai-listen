import type { AdminLoginRequest, AdminLoginResponse, AdminMeResponse } from "@listen/app-shared";
import { adminApiClient } from "./client";

export function adminLogin(payload: AdminLoginRequest) {
  return adminApiClient.post<AdminLoginResponse>("/auth/login", payload);
}

export function adminMe() {
  return adminApiClient.get<AdminMeResponse>("/auth/me");
}
