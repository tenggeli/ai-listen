import type { ActionMeta } from "./common";

export interface UserProfile {
  id: number;
  mobile: string;
  nickname: string;
  avatar: string;
}

export interface SendSMSRequest {
  mobile: string;
  scene?: string;
}

export interface SendSMSResponse extends ActionMeta {
  request: {
    mobile: string;
    scene?: string;
  };
  debugCode: string;
}

export interface LoginBySMSRequest {
  mobile: string;
  code: string;
}

export interface LoginBySMSResponse extends ActionMeta {
  accessToken: string;
  refreshToken: string;
  user: UserProfile;
}

export interface UserMeResponse extends ActionMeta {
  user: UserProfile;
}
