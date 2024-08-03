import React from 'react'
import ReactDOM from 'react-dom/client'
import App from '../src/App.tsx'
import '../src/index.css'

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <h6 style={{ color: "orange" }}>Nested Entry!</h6>
    <App />
  </React.StrictMode >,
)
