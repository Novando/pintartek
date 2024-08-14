import React from 'react'
import { createRoot } from 'react-dom/client'
import { routes } from '@routes/index'
import pkg from '@root/package.json'
import { createBrowserRouter, RouterProvider } from 'react-router-dom'
import NotificationToast from '@arutek/core-app/components/NotificationToast'
import '@styles/index.css'
import '@arutek/core-app/fonticons/dist/aru-icon.scss'

const router = createBrowserRouter(routes)
console.log(`${pkg.displayName} v${pkg.version}`)

createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    <RouterProvider router={router} />
    <NotificationToast />
  </React.StrictMode>,
)
