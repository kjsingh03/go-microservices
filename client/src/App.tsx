import { useRef, useState } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'
import axios from 'axios'

const SERVER_URL = import.meta.env.VITE_SERVER_URL
const AUTH_PORT = import.meta.env.VITE_AUTH_PORT as string;
const MAIL_PORT = import.meta.env.VITE_MAIL_PORT as string;
const LISTENER_PORT = import.meta.env.VITE_LISTENER_PORT as string;
const SERVICE_BROKER_PORT = import.meta.env.VITE_SERVICE_BROKER_PORT as string;

function App() {

  const [loading, setLoading] = useState<boolean>(false)
  const [output, setoutput] = useState<any>("")

  const handleServiceBrokerService = async () => {
    try {
      const res = await axios.get(`${SERVER_URL}:${SERVICE_BROKER_PORT}/broker`)

      setoutput(JSON.stringify(res.data, undefined, 4))
      console.log(res.data)

    } catch (e) {
      e instanceof Error ? console.log(e.message) : console.log(e);
    }
  }
  
  const handleAuthenticationService = async () => {
    try {
      const res = await axios.get(`${SERVER_URL}:${AUTH_PORT}`)

      setoutput(JSON.stringify(res.data, undefined, 4))
      console.log(res.data)

    } catch (e) {
      e instanceof Error ? console.log(e.message) : console.log(e);
    }
  }

  return (
    <div className="w-[80%] mx-auto flex flex-col gap-4">
      <h1 className="mt-5">Test microservices</h1>
      <hr />

      <div className="flex items-center gap-4">
        <button onClick={handleServiceBrokerService} className='rounded-4xl py-2.5 px-3 bg-[#e2e2e2]'>Test Service Broker</button>
        <button onClick={handleAuthenticationService} className='rounded-4xl py-2.5 px-3 bg-[#e2e2e2]'>Test Authentication Service</button>
      </div>

      <div id="output" className="border border-gray-400 p-4">
        <span className="">Output shows here...</span>
      </div>

      <div className="w-full flex gap-4">
        <div className="w-[50%] h-full min-h-64 flex flex-col gap-4">
          <h4 className="">Sent</h4>
          <div className="h-full border border-gray-400 p-4">
            <pre id="payload">
              <span className="">Nothing sent yet...</span>
            </pre>
          </div>
        </div>
        <div className="w-[50%] h-full min-h-64 flex flex-col gap-4">
          <h4 className="">Received</h4>
          <div className="h-full border border-gray-400 p-4 overflow-hidden">
            <pre>
              {
                output ?
                  <>
                    {output}
                  </> :
                  "Nothing received yet..."
              }
            </pre>
          </div>
        </div>
      </div>
    </div>
  )
}

export default App
