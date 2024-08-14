import { RouteObject } from 'react-router-dom'
import Login from '@pages/auth/Login'
import Register from '@pages/auth/Register'
import libMiddleware from '@arutek/frontend-library/src/middleware'
// import Home from '@src/pages/Home'

export const authRoutes:RouteObject[] = [
  {
    path: '/login',
    element: <Login />,
  },
  {
    path: '/register',
    element: <Register />,
  },
  {
    path: '/logout',
    loader: async () => {
      return libMiddleware.logout('/login')
    }
  },
]
