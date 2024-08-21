import {ChangeEvent, Dispatch, SetStateAction} from 'react'

export default (dt: object) => {
  Object.keys(dt).forEach((key) => {
    console.log('coba')
    const doc = document.querySelector(`[name="${key}"]`)
    // @ts-ignore
    if (doc) doc.value = dt[key]
  })
}