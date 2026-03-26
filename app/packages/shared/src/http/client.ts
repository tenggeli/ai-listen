import type { ApiEnvelope, HttpClientConfig, HttpMethod, HttpRequestOptions } from "./types";
import { ApiError } from "./types";

function normalizeBaseURL(baseURL: string): string {
  return baseURL.endsWith("/") ? baseURL.slice(0, -1) : baseURL;
}

function normalizePath(path: string): string {
  return path.startsWith("/") ? path : `/${path}`;
}

export function createHttpClient(config: HttpClientConfig) {
  const baseURL = normalizeBaseURL(config.baseURL);
  const unauthorizedCodes = config.unauthorizedCodes ?? [40100];

  async function request<T>(options: HttpRequestOptions): Promise<T> {
    const accessToken = config.getAccessToken?.();
    const headers: Record<string, string> = {
      "Content-Type": "application/json",
      ...(options.headers ?? {})
    };
    if (accessToken) {
      headers.Authorization = `Bearer ${accessToken}`;
    }

    const response = await new Promise<UniApp.RequestSuccessCallbackResult>((resolve, reject) => {
      uni.request({
        url: `${baseURL}${normalizePath(options.path)}`,
        method: options.method ?? "GET",
        data: options.data,
        header: headers,
        success: resolve,
        fail: reject
      });
    });

    const body = response.data as ApiEnvelope<T> | T;
    if (body && typeof body === "object" && "code" in body && typeof body.code === "number") {
      const envelope = body as ApiEnvelope<T>;
      if (envelope.code === 0) {
        return envelope.data;
      }
      if (unauthorizedCodes.includes(envelope.code)) {
        config.onUnauthorized?.();
      }
      throw new ApiError(envelope.message || "request failed", envelope.code, envelope.traceId, envelope.data);
    }

    if (response.statusCode >= 200 && response.statusCode < 300) {
      return body as T;
    }
    throw new ApiError("request failed", response.statusCode || -1);
  }

  function methodRequest<T>(method: HttpMethod, path: string, data?: unknown) {
    return request<T>({ method, path, data });
  }

  return {
    request,
    get: <T>(path: string) => methodRequest<T>("GET", path),
    post: <T>(path: string, data?: unknown) => methodRequest<T>("POST", path, data),
    put: <T>(path: string, data?: unknown) => methodRequest<T>("PUT", path, data),
    patch: <T>(path: string, data?: unknown) => methodRequest<T>("PATCH", path, data),
    del: <T>(path: string, data?: unknown) => methodRequest<T>("DELETE", path, data)
  };
}
