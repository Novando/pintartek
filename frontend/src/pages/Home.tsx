import { useState, useEffect, useRef, ChangeEvent } from 'react'
import libDate from '@arutek/core-app/libraries/date'
import {Link, useNavigate} from 'react-router-dom'
import vault, {VaultResponseType} from '@factories/vault'
import notify from '@arutek/core-app/helpers/notification'
import NewVaultModal from '@src/components/modal/NewVaultModal'
import {closeModal, showModal} from '@src/utils/modal'

const Home = () => {
  const [vaults, setVaults] = useState<VaultResponseType[]>([])
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
          <div onClick={() => showModal('newVaultModal')} className="text-center w-[160px] h-[160px] bg-neutral-600 cursor-pointer">
            <p>Add new vault</p>
          </div>
          {vaults.map((vault, key) => (
            <div
              key={key}
              onClick={() => navigate(`/vault/${vault.id}`)}
              className="text-center w-[160px] h-[160px] bg-neutral-600 cursor-pointer">
              <p>{vault.name}</p>
            </div>
          ))}
        </section>
      </section>
      <NewVaultModal onCreateVault={() => closeModal('newVaultModal')} modalId="newVaultModal"/>
    </main>
  )
}

export default Home
