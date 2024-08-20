import { useState, useEffect, useRef, ChangeEvent } from 'react'
import libDate from '@arutek/core-app/libraries/date'
import callModal from '@src/utils/call-modal'
import {Link, useNavigate, useParams} from 'react-router-dom'
import vault, {CredentialType} from '@factories/vault'
import CredentialModal from '@src/components/modal/CredentialModal'
import {useNotification} from '@src/components/NotificationToast'
import {TrashCan} from '@src/components/svg/TrashCan'
import {Copy} from '@src/components/svg/Copy'

const Vault = () => {
  const [credentials, setCredentials] = useState<(CredentialType & {id: string})[]>([])
  const {addNoty} = useNotification()
  const {vaultId} = useParams()

  useEffect(() => {
    init()
  }, [])

  const init = async () => {
    try {
      const res = await vault.getOne(vaultId || '')
      const json = JSON.parse(atob(res.data))
      const ids = Object.keys(json)
      const newCredentials: (CredentialType & {id: string})[] = []
      for (const id of ids) {
        newCredentials.push({
          id,
          name: json[id].name,
          credential: json[id].credential,
          url: json[id].url,
          note: json[id].note,
          password: json[id].password,
        })
      }
      setCredentials(newCredentials)
    } catch (e: any) {
      addNoty(e.message, 'error')
    }
  }

  const clipboardCopy = (val: string) => {
    navigator.clipboard.writeText(val)
    addNoty('Password has been copied to clipboard', 'success', 'Copied!')
  }

  return (
    <main>
      <section className="py-4 px-8 bg-sky-800 text-white flex justify-end">
        <Link to="/logout">Logout</Link>
      </section>
      <section className="mx-auto max-w-7xl">
        <section className="my-8">
          <table className="w-full">
            <thead>
            <tr>
              <th>Name</th>
              <th>Credential</th>
              <th>Password</th>
              <th>Created At</th>
              <th>Action</th>
            </tr>
            </thead>
            <tbody>
              {credentials.map((credential) => (
                <tr>
                  <td>{credential.name}</td>
                  <td onClick={() => clipboardCopy(credential.credential)} className="cursor-pointer">
                    <div className="flex justify-center gap-4">
                      <span>********</span>
                      <div className="w-4">
                        <Copy/>
                      </div>
                    </div>
                  </td>
                  <td>{credential.password}</td>
                  <td>{libDate.isoToDate1('2024-08-08T08:10:00Z')}</td>
                  <td>
                    <div className="flex gap-4">
                      <button onClick={() => callModal()}>
                        <div className="w-4">
                          <TrashCan />
                        </div>
                      </button>
                      <p>V</p>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </section>
      </section>
      <CredentialModal/>
    </main>
  )
}

export default Vault
