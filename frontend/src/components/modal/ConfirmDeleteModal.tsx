import vault from '@factories/vault'
import notify from '@arutek/core-app/helpers/notification'
import {useParams} from 'react-router-dom'

type ConfirmDeleteModalProps = {
  modalId: string
  credentialId: string
  name: string
  onSuccess: () => void
}

const ConfirmDeleteModal = (props: ConfirmDeleteModalProps) => {
  const {vaultId} = useParams()

  const deleteCredential = async () => {
    try {
      if (!vaultId) throw ({message: 'Vault ID not found'})
      await vault.deleteCredential(vaultId, props.credentialId)
      props.onSuccess()
    } catch (e: any) {
      notify.notifyError(e.message)
    }
  }

  return (
    <dialog id={props.modalId} className="modal">
      <div className="modal-box">
        <form method="dialog">
          {/* if there is a button in form, it will close the modal */}
          <button className="btn btn-sm btn-circle btn-ghost absolute right-2 top-2">✕</button>
        </form>
        <h3 className="font-bold text-lg">Confirm Credential Deletion</h3>
        <p>Are you sure want to delete <span className="font-bold">{props.name}</span> permanently?</p>
        <div className="flex gap-4 justify-center">
          <form method="dialog">
            <button className="bg-neutral-400 text-black rounded px-4 py-1">Cancel</button>
          </form>
          <button onClick={deleteCredential} className="bg-sky-400 text-black rounded px-4 py-1">Create</button>
        </div>
      </div>
      <form method="dialog" className="modal-backdrop">
        <button className="cursor-default">close</button>
      </form>
    </dialog>
  )
}

export default ConfirmDeleteModal