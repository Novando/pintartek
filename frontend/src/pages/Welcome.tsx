import { useState, useEffect } from 'react'
import {useNavigate, useSearchParams} from 'react-router-dom'

const Welcome = () => {
  const navigate = useNavigate()
  const [searchParams, setSearchParams] = useSearchParams()
  const [privateKey, setPrivateKey] = useState('')

  useEffect(() => {
    init()
  }, [])

  const init = () => {
    setPrivateKey(searchParams.get('key') || '')
  }

  return (
    <main>
      <section className="mx-auto max-w-7xl">
        <section className="my-8">
          <h1 className="text-center text-xl font-bold">Welcome</h1>
        </section>
        <section className="max-w-3xl mx-auto">
          <p className="mb-4">
            Make sure to save this private key in case you need to reset your password.<br/>
            <span className="font-bold">This key will never be found anywhere again, </span>
            save it on very secure and private environment.
          </p>
          <div className="bg-neutral-600 w-full py-1 px-2 mb-10">
            <pre className="whitespace-pre-wrap break-all">{privateKey}</pre>
          </div>
          <button
            onClick={() => navigate('/login', {replace: true})}
            className="bg-sky-400 text-black rounded px-4 py-1">
            Login
          </button>
        </section>
      </section>
    </main>
  )
}

export default Welcome
