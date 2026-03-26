import type { HealthStatus } from "@listen/app-shared";
import { userApiClient } from "./client";

export function getHealth() {
  return userApiClient.get<HealthStatus>("/health");
}
