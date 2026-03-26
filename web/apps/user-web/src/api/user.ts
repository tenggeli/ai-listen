import { request } from "./client";

export interface HealthInfo {
  service: string;
  status: string;
}

export interface User {
  id: number;
  mobile: string;
  nickname: string;
  avatar: string;
  gender: number;
  birthday: string;
  cityCode: string;
  createdAt: string;
  updatedAt: string;
}

export interface Order {
  id: number;
  orderNo: string;
  userId: number;
  providerId: number;
  serviceItemId: number;
  sceneText: string;
  cityCode: string;
  addressText: string;
  plannedStartAt: string;
  plannedDuration: number;
  status: number;
  payAmount: number;
  createdAt: string;
}

interface SendSMSResponse {
  debugCode?: string;
}

interface LoginResponse {
  accessToken: string;
  refreshToken: string;
  user: User;
}

interface MeResponse {
  user: User;
}

interface OrdersResponse {
  list: Order[];
}

export function getHealth() {
  return request<HealthInfo>("/health");
}

export function sendSMS(mobile: string) {
  return request<SendSMSResponse>("/auth/sms/send", {
    method: "POST",
    body: { mobile, scene: "login" }
  });
}

export function loginBySMS(mobile: string, code: string) {
  return request<LoginResponse>("/auth/login/sms", {
    method: "POST",
    body: { mobile, code }
  });
}

export function getMe() {
  return request<MeResponse>("/users/me", {
    auth: true
  });
}

export function getMyOrders() {
  return request<OrdersResponse>("/users/me/orders", {
    auth: true
  });
}
