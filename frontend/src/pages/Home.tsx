import { useState, useEffect, useRef, ChangeEvent } from 'react'
import libDate from '@arutek/core-app/libraries/date'
import Modal from '@src/components/Modal'
import callModal from '@src/utils/call-modal'
import {Link} from 'react-router-dom'

const Home = () => {
  let letOfficialSearch:string
  
  useEffect(() => {
    init()
  }, [])

  const init = () => {

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
              <th>Password</th>
              <th>Created At</th>
              <th>Action</th>
            </tr>
            </thead>
            <tbody>
            <tr>
              <td>pass 1</td>
              <td>********</td>
              <td>{libDate.isoToDate1('2024-08-08T08:10:00Z')}</td>
              <td>
                <div className="flex gap-4">
                  <button
                    onClick={() => callModal()}>
                    D
                  </button>
                  <p>V</p>
                </div>
              </td>
            </tr>
            </tbody>
          </table>
        </section>
      </section>
      <Modal />
    </main>
  )
}

export default Home
