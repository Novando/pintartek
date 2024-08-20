import handleInput from '@src/utils/handle-input'
import {useState} from 'react'
import vault, {CredentialType} from '@factories/vault'
import notify from '@arutek/core-app/helpers/notification'
import {useParams} from 'react-router-dom'


const CredentialModal = () => {
  const {vaultId} = useParams()
  const [credential, setCredential] = useState<CredentialType>({
    name: '',
    credential: '',
    password: '',
    url: '',
    note: '',
  })

  const create = async () => {
    try {
      await vault.createCredential(vaultId || '', credential)
    } catch (e: any) {
      notify.notifyError(e.message)
    }
  }

  return (
    <dialog id="modal" className="modal">
      <div className="modal-box">
        <form method="dialog">
          {/* if there is a button in form, it will close the modal */}
          <button className="btn btn-sm btn-circle btn-ghost absolute right-2 top-2">✕</button>
        </form>
        <h3 className="font-bold text-lg">Credential Wizard</h3>
        <form onSubmit={(e) => {e.preventDefault()}} className="py-4 mb-6">
          <label>
            <p className="mb-1">Credential Name</p>
            <input
              className="text-black bg-white py-1 px-2 rounded w-2/3 mb-2"
              type="text"
              placeholder="Your credential name"
              name="name"
              onChange={(e) => handleInput(e, setCredential)}/>
          </label>
          <label>
            <p className="mb-1">Credential</p>
            <input
              className="text-black bg-white py-1 px-2 rounded w-2/3 mb-2"
              type="text"
              placeholder="Your credential"
              name="credential"
              onChange={(e) => handleInput(e, setCredential)}/>
          </label>
          <label>
            <p className="mb-1">Password</p>
            <input
              className="text-black bg-white py-1 px-2 rounded w-2/3 mb-2"
              type="text"
              placeholder="Your password"
              name="password"
              onChange={(e) => handleInput(e, setCredential)}/>
          </label>
          <label>
            <p className="mb-1">URL</p>
            <input
              className="text-black bg-white py-1 px-2 rounded w-2/3 mb-2"
              type="text"
              placeholder="Credential URL"
              name="url"
              onChange={(e) => handleInput(e, setCredential)}/>
          </label>
          <label>
            <p className="mb-1">Note</p>
            <textarea
              className="text-black bg-white py-1 px-2 rounded w-2/3 mb-2 resize-none"
              rows={3}
              name="note"
              onChange={(e) => handleInput(e, setCredential)}/>
          </label>
          <button onClick={create} className="hidden">Login</button>
        </form>
        <div>
          <button onClick={create} className="bg-sky-400 text-black rounded px-4 py-1">Create</button>
        </div>
      </div>
      <form method="dialog" className="modal-backdrop">
        <button className="cursor-default">close</button>
      </form>
    </dialog>
  )
}

export default CredentialModal