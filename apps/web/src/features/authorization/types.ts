export namespace AuthorizationServices {
  export interface Token {
    accessToken: string
  }

  export type LoginInfo = Record<'username' | 'password', string>
  export type LoginInfoError = LoginInfo

  export interface UserInfo {
    uid: string
    username: string
    nickname: string
    userGroup: string
    icon: string
  }

  export interface RegisterInfo {
    username: string
    password: string
    confirmPassword: string
  }

  export type RegisterInfoError = RegisterInfo
}
