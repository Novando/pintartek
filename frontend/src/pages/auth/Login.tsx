import {Link, useNavigate} from 'react-router-dom'
import userFactory, {LoginParamType} from '@factories/user'
import helpCookie from '@arutek/core-app/helpers/cookie'
import {useState} from 'react'
import handleInput from '@src/utils/handle-input'
import {useNotification} from '@src/components/NotificationToast'

const Login = () => {
  const { addNoty } = useNotification()
  const [loginPayload, setLoginPayload] = useState<LoginParamType>({
    email: '',
    password: '',
  })

  const navigate = useNavigate()

  const login = async () => {
    try {
      const res = await userFactory.login(loginPayload)
      helpCookie.setAuthCookie(res.data.accessToken, 30)
      helpCookie.setCookie('userData', '{"roleId":0}', 60*24*365)
      setTimeout(() => navigate('/', {replace: true}), 250)
    } catch (e: any) {
      addNoty(e.message, 'error')
    }
  }

  return (
    <main>
      <section className="text-white mx-auto mt-40 py-8 px-12 bg-sky-800 rounded-lg border shadow w-[480px]">
        <h1 className="text-xl mb-6">Login</h1>
        <form onSubmit={(e) => {e.preventDefault()}} className="mb-6">
          <label>
            <p className="mb-1">Email</p>
            <input
              className="text-black bg-white py-1 px-2 rounded w-2/3 mb-2"
              type="email"
              placeholder="Your email address"
              name="email"
              onChange={(e) => handleInput(e, setLoginPayload)}/>
          </label>
          <label>
            <p className="mb-1">Password</p>
            <input
              className="text-black bg-white py-1 px-2 rounded w-2/3"
              type="password"
              placeholder="Your password"
              name="password"
              onChange={(e) => handleInput(e, setLoginPayload)}/>
          </label>
          <button onClick={login} className="hidden">Login</button>
        </form>
        <div>
          <div className="mb-2">
            <p>Did not have an account?</p>
            <Link to="/register" className="text-sky-400">Create an account</Link>
          </div>
          <button onClick={login} className="bg-sky-400 text-black rounded px-4 py-1">Login</button>
        </div>
      </section>
    </main>
  )
}

export default Login