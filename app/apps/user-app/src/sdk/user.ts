import type { UserMeResponse } from "@listen/app-shared";
import { userApiClient } from "./client";

export function getMe() {
  return userApiClient.get<UserMeResponse>("/users/me");
}
