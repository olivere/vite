import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App.tsx'
import './index.css'

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <h6 style={{ color: "orange" }}>Index Entry!</h6>
    <App />
  </React.StrictMode>,
)
