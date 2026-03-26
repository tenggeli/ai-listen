export type HttpMethod = "GET" | "POST" | "PUT" | "PATCH" | "DELETE";

export interface ApiEnvelope<T> {
  code: number;
  message: string;
  data: T;
  traceId: string;
}

export interface HttpRequestOptions {
  method?: HttpMethod;
  path: string;
  data?: unknown;
  headers?: Record<string, string>;
}

export interface HttpClientConfig {
  baseURL: string;
  getAccessToken?: () => string;
  onUnauthorized?: () => void;
  unauthorizedCodes?: number[];
}

export class ApiError<T = unknown> extends Error {
  code: number;
  traceId?: string;
  data?: T;

  constructor(message: string, code: number, traceId?: string, data?: T) {
    super(message);
    this.name = "ApiError";
    this.code = code;
    this.traceId = traceId;
    this.data = data;
  }
}
