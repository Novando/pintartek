import {createContext, useContext, useState} from 'react'
import {debounce} from 'lodash'

type NotificationPropsType = {
  message: string
  type?: string
  title?: string
}

interface NotificationToastProps {
  addNoty: (message: string, type?: string, title?: string) => void
}

const NotificationToast = createContext<NotificationToastProps | undefined>(undefined)
export const useNotification = (): NotificationToastProps => {
  const context = useContext(NotificationToast)
  if (!context) {
    throw new Error('useToast must be used within a ToastProvider');
  }
  return context
}
export const NotificationProvider = (props: {children: JSX.Element}) => {
  const [notys, setNotys] = useState<NotificationPropsType[]>([])

  const addNoty = (message: string, type = 'success', title?: string) => {
    const newNotys = [...notys, { message, type, title }]
    if (newNotys.length > 3) newNotys.shift()
    setNotys(newNotys)
    debounceReducer()
  }

  // const reduceNoty = (params = [...notys]) => {
  //   if (params.length > 0) setTimeout(() => {
  //     params.shift()
  //     setNotys(params)
  //     reduceNoty(params)
  //   }, 500)
  // }

  const debounceReducer = debounce(() => setNotys([]), 2000)

  return (
    <NotificationToast.Provider value={{addNoty}}>
      {props.children}
      <div className="fixed h-full bottom-4 right-4">
          <div className="toast">
            {notys.map((noty, key) => (
              <div key={key} className="alert bg-red-400 text-black">
                <span>{noty.message}</span>
              </div>
            ))}
        </div>
      </div>
    </NotificationToast.Provider>
  )
}