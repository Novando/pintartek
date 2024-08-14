import libFetch from '@arutek/core-app/libraries/fetch'

const apiUrl = `${import.meta.env.VITE_API_URL}/reason`

export default {
  getAll ():Promise<responseType> {
    return libFetch.getData(`${apiUrl}`)
  },
}