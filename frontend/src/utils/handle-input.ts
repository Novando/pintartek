import {ChangeEvent, Dispatch, SetStateAction} from 'react'

export default (e: ChangeEvent<HTMLInputElement|HTMLTextAreaElement>, setter: Dispatch<SetStateAction<any>>) => {
  setter((prevState: any) => {
      return {...prevState, [e.target.name]: e.target.value}
  })
}