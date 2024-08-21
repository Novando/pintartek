export const showModal = (id = 'modal') => {
  // @ts-ignore
  document.getElementById(id).showModal()
}

export const closeModal = (id = 'modal') => {
  // @ts-ignore
  document.getElementById(id).close()
}