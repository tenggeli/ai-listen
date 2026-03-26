import type { DashboardOverviewResponse } from "@listen/app-shared";
import { adminApiClient } from "./client";

export function dashboardOverview() {
  return adminApiClient.get<DashboardOverviewResponse>("/dashboard/overview");
}
