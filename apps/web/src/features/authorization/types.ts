export namespace AuthorizationServices {
  export interface Token {
    accessToken: string;
  }

  export type LoginInfo = Record<"username" | "password", string>;
  export type LoginInfoError = LoginInfo;

  export interface UserInfo {
    uid: number;
    username: string;
    nickname: string;
    userGroup: string;
    mail: string;
    icon: string;
  }

  export interface RegisterInfo {
    username: string;
    password: string;
    confirmPassword: string;
  }

  export type RegisterInfoError = RegisterInfo;
}
