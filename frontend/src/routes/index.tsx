import { RouteObject } from 'react-router-dom'
import Home from '@pages/Home'
import libMiddleware from '@arutek/frontend-library/src/middleware'
import {authRoutes} from '@routes/auth'

export const routes:RouteObject[] = [
  ...authRoutes,
  {
    path: '/',
    element: <Home />,
    // loader: async () => {
    //   return libMiddleware.isAuthenticated('/login')
    // }
  },
]