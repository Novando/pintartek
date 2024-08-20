import libFetch from '@arutek/core-app/libraries/fetch'

type CreatePayloadType = {
  name: string
  credential: CredentialType
}
export type CredentialType = {
  name: string
  password: string
  credential: string
  url: string
  note: string
}
export type VaultResponseType = {
  id: string
  name: string
  createdAt: string
}

const apiUrl = `${import.meta.env.VITE_API_URL}/vault`

export default {
  getAll ():Promise<ResponseAPIType & {data: VaultResponseType[]}> {
    return libFetch.getDataLogged(apiUrl)
  },
  getOne (id: string):Promise<ResponseAPIType> {
    return libFetch.getDataLogged(`${apiUrl}/${id}`)
  },
  create (payload: CreatePayloadType):Promise<ResponseAPIType> {
    return libFetch.postDataLogged(`${apiUrl}`, payload)
  },
  createCredential (vaultId: string, payload: CreatePayloadType):Promise<ResponseAPIType> {
    return libFetch.postDataLogged(`${apiUrl}/${vaultId}`, payload)
  },
  update (vaultId: string, payload: { name: string }):Promise<ResponseAPIType> {
    return libFetch.putDataLogged(`${apiUrl}/${vaultId}`, payload)
  },
  updateCredential (vaultId: string, credentialId: string, payload: CreatePayloadType):Promise<ResponseAPIType> {
    return libFetch.putDataLogged(`${apiUrl}/${vaultId}/${credentialId}`, payload)
  },
  delete (vaultId: string):Promise<ResponseAPIType> {
    return libFetch.delDataLogged(`${apiUrl}/${vaultId}`)
  },
  deleteCredential (vaultId: string, credentialId: string):Promise<ResponseAPIType> {
    return libFetch.delDataLogged(`${apiUrl}/${vaultId}/${credentialId}`)
  },
}