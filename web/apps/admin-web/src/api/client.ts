export interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
  traceId: string;
}

const TOKEN_KEY = "listen.admin.token";
const BASE_URL = import.meta.env.VITE_ADMIN_API_BASE_URL ?? "/api/v1/admin";

export function getAdminAccessToken() {
  return localStorage.getItem(TOKEN_KEY) ?? "";
}

export function setAdminAccessToken(token: string) {
  localStorage.setItem(TOKEN_KEY, token);
}

export function clearAdminAccessToken() {
  localStorage.removeItem(TOKEN_KEY);
}

interface RequestOptions {
  method?: "GET" | "POST" | "PUT" | "DELETE";
  body?: unknown;
  auth?: boolean;
}

export async function request<T>(path: string, options: RequestOptions = {}) {
  const method = options.method ?? "GET";
  const headers: Record<string, string> = {
    "Content-Type": "application/json"
  };

  if (options.auth) {
    const token = getAdminAccessToken();
    if (token) {
      headers.Authorization = `Bearer ${token}`;
    }
  }

  const response = await fetch(`${BASE_URL}${path}`, {
    method,
    headers,
    body: options.body == null ? undefined : JSON.stringify(options.body)
  });
  const payload = (await response.json()) as ApiResponse<T>;

  if (!response.ok || payload.code !== 0) {
    if (response.status === 401) {
      clearAdminAccessToken();
    }
    throw new Error(payload.message || "请求失败");
  }

  return payload.data;
}
