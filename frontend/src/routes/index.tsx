import { RouteObject } from 'react-router-dom'
import Home from '@pages/Home'
import libMiddleware from '@arutek/frontend-library/src/middleware'
import {authRoutes} from '@routes/auth'
import Welcome from '@pages/Welcome'
import Vault from '@pages/Vault'

export const routes:RouteObject[] = [
  ...authRoutes,
  {
    path: '/',
    element: <Home />,
    loader: async () => {
      return libMiddleware.isAuthenticated('/login', 0)
    },
  },
  {
    path: '/vault/:vaultId',
    element: <Vault />,
    loader: async () => {
      return libMiddleware.isAuthenticated('/login', 0)
    },
  },
  {
    path: '/welcome',
    element: <Welcome />,
  },
]