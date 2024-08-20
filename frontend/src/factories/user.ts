import libFetch from '@arutek/core-app/libraries/fetch'

export type LoginParamType = {
  email: string
  password: string
}
export type RegisterParamType = LoginParamType & {
  fullName: string
  confirmPassword: string
}

const apiUrl = `${import.meta.env.VITE_API_URL}/user`

export default {
  register (payload: RegisterParamType):Promise<ResponseAPIType & {data: { privateKey: string }}> {
    return libFetch.postData(`${apiUrl}/register`, payload)
  },
  login (payload: LoginParamType):Promise<ResponseAPIType & {data: { accessToken: string }}> {
    return libFetch.postData(`${apiUrl}/login`, payload)
  },
  logout ():Promise<ResponseAPIType> {
    return libFetch.getData(`${apiUrl}/logout`)
  },
}