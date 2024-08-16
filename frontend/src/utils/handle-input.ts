import {ChangeEvent, Dispatch, SetStateAction} from 'react'

export default (e: ChangeEvent<HTMLInputElement>, setter: Dispatch<SetStateAction<any>>) => {
  setter((prevState: any) => {
      return {...prevState, [e.target.name]: e.target.value}
  })
}