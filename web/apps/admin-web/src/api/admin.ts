import { request } from "./client";

export interface AdminUser {
  id: number;
  username: string;
  nickname: string;
  status: number;
  roles: string[];
  createdAt: string;
  updatedAt: string;
}

export interface PaginationQuery {
  page: number;
  pageSize: number;
}

export interface AdminLoginResponse {
  accessToken: string;
  adminUser: AdminUser;
  permissions: string[];
}

export interface AdminMeResponse {
  adminUser: AdminUser;
  permissions: string[];
}

export interface AdminActionMeta {
  module: string;
  action: string;
}

export interface ListPostsResponse extends AdminActionMeta {
  query: PaginationQuery;
}

export interface ListAudioResponse extends AdminActionMeta {
  query: PaginationQuery;
}

export interface HidePostResponse extends AdminActionMeta {
  postId: string;
}

export interface OffShelfAudioResponse extends AdminActionMeta {
  audioId: string;
}

export interface ListComplaintsResponse extends AdminActionMeta {
  query: PaginationQuery;
}

export interface ComplaintDetailResponse extends AdminActionMeta {
  complaintId: string;
}

export interface ResolveComplaintResponse extends AdminActionMeta {
  complaintId: string;
}

export interface ListRiskEventsResponse extends AdminActionMeta {
  query: PaginationQuery;
}

export function adminLogin(username: string, password: string) {
  return request<AdminLoginResponse>("/auth/login", {
    method: "POST",
    body: { username, password }
  });
}

export function adminMe() {
  return request<AdminMeResponse>("/auth/me", {
    auth: true
  });
}

export function adminLogout() {
  return request<AdminActionMeta>("/auth/logout", {
    method: "POST",
    auth: true
  });
}

export function listPosts(page = 1, pageSize = 20) {
  return request<ListPostsResponse>(`/posts?page=${page}&pageSize=${pageSize}`, {
    auth: true
  });
}

export function hidePost(postId: number) {
  return request<HidePostResponse>(`/posts/${postId}/hide`, {
    method: "POST",
    auth: true
  });
}

export function listAudio(page = 1, pageSize = 20) {
  return request<ListAudioResponse>(`/audio?page=${page}&pageSize=${pageSize}`, {
    auth: true
  });
}

export function offShelfAudio(audioId: number) {
  return request<OffShelfAudioResponse>(`/audio/${audioId}/off-shelf`, {
    method: "POST",
    auth: true
  });
}

export function listComplaints(page = 1, pageSize = 20) {
  return request<ListComplaintsResponse>(`/complaints?page=${page}&pageSize=${pageSize}`, {
    auth: true
  });
}

export function complaintDetail(complaintId: number) {
  return request<ComplaintDetailResponse>(`/complaints/${complaintId}`, {
    auth: true
  });
}

export function resolveComplaint(complaintId: number) {
  return request<ResolveComplaintResponse>(`/complaints/${complaintId}/resolve`, {
    method: "POST",
    auth: true
  });
}

export function listRiskEvents(page = 1, pageSize = 20) {
  return request<ListRiskEventsResponse>(`/risk-events?page=${page}&pageSize=${pageSize}`, {
    auth: true
  });
}
