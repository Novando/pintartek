import { useState, useEffect, useRef, ChangeEvent } from 'react'
import libDate from '@arutek/core-app/libraries/date'
import Modal from '@src/components/Modal'
import callModal from '@src/utils/call-modal'
import {Link, useNavigate} from 'react-router-dom'
import vault from '@factories/vault'
import notify from '@arutek/core-app/helpers/notification'
import NewVaultModal from '@src/components/modal/NewVaultModal'

const Home = () => {
  const [vaults, setVaults] = useState([])
  const navigate = useNavigate()

  useEffect(() => {
    init()
  }, [])

  const init = async () => {
    try {
      const res = await vault.getAll()
      setVaults(res.data)
    } catch (e: any) {
      notify.notifyError(e.message)
    }

  }

  return (
    <main>
      <section className="py-4 px-8 bg-sky-800 text-white flex justify-end">
        <Link to="/logout">Logout</Link>
      </section>
      <section className="mx-auto max-w-7xl">
        <section className="my-8">
          <h1 className="text-xl font-bold text-center">Vault List</h1>
        </section>
        <section className="flex gap-6">
          <div onClick={callModal} className="text-center w-[160px] h-[160px] bg-neutral-600">
            <p>Add new vault</p>
          </div>
          {vaults.map((vault) => (
            <div onClick={() => navigate(`/vault/${vault.id}`)} className="text-center w-[160px] h-[160px] bg-neutral-600">
              <p>{vault.name}</p>
            </div>
          ))}
        </section>
      </section>
      <NewVaultModal/>
    </main>
  )
}

export default Home
