import {Link, useNavigate} from 'react-router-dom'
import userFactory, {RegisterParamType} from '@factories/user'
import helpCookie from '@arutek/core-app/helpers/cookie'
import notify from '@arutek/core-app/helpers/notification'
import {useState} from 'react'
import handleInput from '@src/utils/handle-input'


const Register = () => {
  const [registerForm, setRegisterForm] = useState<RegisterParamType>({
    fullName: '',
    email: '',
    password: '',
    confirmPassword: '',
  })

  const navigate = useNavigate()

  const register = async () => {
    try {
      const res = await userFactory.register(registerForm)
      navigate(`/welcome?key=${res.data.privateKey}`, {replace: true})
    } catch (e: any) {
      notify.notifyError(e.message)
    }
  }
  return (
    <main>
      <section className="text-white mx-auto mt-40 py-8 px-12 bg-sky-800 rounded-lg border shadow w-[480px]">
        <h1 className="text-xl mb-6">Register</h1>
        <form onSubmit={(e) => e.preventDefault()} className="mb-6">
          <label>
            <p className="mb-1">Full Name</p>
            <input
              type="text"
              name="fullName"
              onChange={(e) => handleInput(e, setRegisterForm)}
              className="text-black bg-white py-1 px-2 rounded w-2/3 mb-2"
              placeholder="Your full name"/>
          </label>
          <label>
            <p className="mb-1">Email</p>
            <input
              type="email"
              name="email"
              onChange={(e) => handleInput(e, setRegisterForm)}
              className="text-black bg-white py-1 px-2 rounded w-2/3 mb-2"
              placeholder="Your email address"/>
          </label>
          <label>
            <p className="mb-1">Password</p>
            <input
              type="password"
              name="password"
              onChange={(e) => handleInput(e, setRegisterForm)}
              className="text-black bg-white py-1 px-2 rounded w-2/3 mb-2"
              placeholder="Your password"/>
          </label>
          <label>
            <p className="mb-1">Confirm Password</p>
            <input
              type="password"
              name="confirmPassword"
              onChange={(e) => handleInput(e, setRegisterForm)}
              className="text-black bg-white py-1 px-2 rounded w-2/3"
              placeholder="Retype Your password"/>
          </label>
          <button onClick={register} className="hidden">Login</button>
        </form>
        <div>
          <div className="mb-2">
            <p>Already register?</p>
            <Link to="/login" className="text-sky-400">Login now</Link>
          </div>
          <button onClick={register} className="bg-sky-400 text-black rounded px-4 py-1">Register</button>
        </div>
      </section>
    </main>
  )
}

export default Register