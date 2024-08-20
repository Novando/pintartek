import React from 'react'
import { createRoot } from 'react-dom/client'
import { routes } from '@routes/index'
import pkg from '@root/package.json'
import { createBrowserRouter, RouterProvider } from 'react-router-dom'
import '@styles/index.css'
import '@arutek/core-app/fonticons/dist/aru-icon.scss'
import {NotificationProvider} from '@src/components/NotificationToast'

const router = createBrowserRouter(routes)
console.log(`${pkg.displayName} v${pkg.version}`)

createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    <NotificationProvider>
      <RouterProvider router={router} />
    </NotificationProvider>
  </React.StrictMode>,
)
