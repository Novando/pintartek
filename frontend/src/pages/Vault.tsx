import { useState, useEffect, useRef, ChangeEvent } from 'react'
import libDate from '@arutek/core-app/libraries/date'
import {Link, useNavigate, useParams} from 'react-router-dom'
import vault, {CredentialType} from '@factories/vault'
import CredentialModal from '@src/components/modal/CredentialModal'
import {useNotification} from '@src/components/NotificationToast'
import {TrashCan} from '@src/components/svg/TrashCan'
import {Copy} from '@src/components/svg/Copy'
import ConfirmDeleteModal from '@src/components/modal/ConfirmDeleteModal'
import {Eye} from '@src/components/svg/Eye'
import {closeModal, showModal} from '@src/utils/modal'

type VaultType = CredentialType & {id: string}

const Vault = () => {
  const [credentials, setCredentials] = useState<VaultType[]>([])
  const [selectedCredential, setSelectedCredential] = useState<VaultType|undefined>(undefined)
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
  const refreshVault = (modalId: string) => {
    init()
    closeModal(modalId)
  }
  const clipboardCopy = (val: string) => {
    navigator.clipboard.writeText(val)
    addNoty('Password has been copied to clipboard', 'success', 'Copied!')
  }

  const confirmDelete = (credential: VaultType) => {
    setSelectedCredential(credential)
    showModal('confirmDeleteModal')
  }


  const credentialDetail = (credential: VaultType) => {
    setSelectedCredential(credential)
    showModal('updateCredentialModal')
  }

  return (
    <main>
      <section className="py-4 px-8 bg-sky-800 text-white flex justify-end">
        <Link to="/logout">Logout</Link>
      </section>
      <section className="mx-auto max-w-7xl">
        <section className="my-8">
          <div className="mb-8">
            <button onClick={() => showModal('createCredentialModal')} className="bg-sky-400 text-black rounded px-4 py-1">Create</button>
          </div>
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
              {credentials.map((credential, key) => (
                <tr key={key}>
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
                      <button onClick={() => confirmDelete(credential)}>
                        <div className="w-4">
                          <TrashCan/>
                        </div>
                      </button>
                      <button onClick={() => credentialDetail(credential)}>
                        <div className="w-4">
                          <Eye/>
                        </div>
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </section>
      </section>
      <ConfirmDeleteModal
        onSuccess={() => refreshVault('confirmDeleteModal')}
        modalId={'confirmDeleteModal'}
        name={selectedCredential?.name || ''}
        credentialId={selectedCredential?.id || ''}/>
      <CredentialModal
        onSuccess={() => refreshVault('updateCredentialModal')}
        modalId={'updateCredentialModal'}
        type="update"
        credential={selectedCredential} />
      <CredentialModal
        onSuccess={() => refreshVault('createCredentialModal')}
        modalId={'createCredentialModal'}
        type="create" />
    </main>
  )
}

export default Vault
