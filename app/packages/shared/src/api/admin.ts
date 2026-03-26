import type { ActionMeta } from "./common";

export interface AdminUser {
  id: number;
  username: string;
  nickname: string;
  roles: string[];
}

export interface AdminLoginRequest {
  username: string;
  password: string;
}

export interface AdminLoginResponse extends ActionMeta {
  accessToken: string;
  adminUser: AdminUser;
  permissions: string[];
}

export interface AdminMeResponse extends ActionMeta {
  adminUser: AdminUser;
  permissions: string[];
}

export interface DashboardOverviewResponse extends ActionMeta {}
