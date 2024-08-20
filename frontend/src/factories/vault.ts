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

const apiUrl = `${import.meta.env.VITE_API_URL}/vault`

export default {
  getAll ():Promise<ResponseType> {
    return libFetch.getDataLogged(apiUrl)
  },
  getOne (id: string):Promise<ResponseType> {
    return libFetch.getDataLogged(`${apiUrl}/${id}`)
  },
  create (payload: CreatePayloadType):Promise<ResponseType> {
    return libFetch.postDataLogged(`${apiUrl}`, payload)
  },
  createCredential (vaultId: string, payload: CreatePayloadType):Promise<ResponseType> {
    return libFetch.postDataLogged(`${apiUrl}/${vaultId}`, payload)
  },
  update (vaultId: string, payload: { name: string }):Promise<ResponseType> {
    return libFetch.putDataLogged(`${apiUrl}/${vaultId}`, payload)
  },
  updateCredential (vaultId: string, credentialId: string, payload: CreatePayloadType):Promise<ResponseType> {
    return libFetch.putDataLogged(`${apiUrl}/${vaultId}/${credentialId}`, payload)
  },
  delete (vaultId: string):Promise<ResponseType> {
    return libFetch.delDataLogged(`${apiUrl}/${vaultId}`)
  },
  deleteCredential (vaultId: string, credentialId: string):Promise<ResponseType> {
    return libFetch.delDataLogged(`${apiUrl}/${vaultId}/${credentialId}`)
  },
}