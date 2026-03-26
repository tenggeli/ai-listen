import type {
  LoginBySMSRequest,
  LoginBySMSResponse,
  SendSMSRequest,
  SendSMSResponse
} from "@listen/app-shared";
import { userApiClient } from "./client";

export function sendSMS(payload: SendSMSRequest) {
  return userApiClient.post<SendSMSResponse>("/auth/sms/send", payload);
}

export function loginBySMS(payload: LoginBySMSRequest) {
  return userApiClient.post<LoginBySMSResponse>("/auth/login/sms", payload);
}
