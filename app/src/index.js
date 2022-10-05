import React from 'react'
import { createRoot } from 'react-dom/client';
import App from './App'

const rootElement = document.getElementById('root');
const root = createRoot(rootElement); // createRoot(rootElement!) if you use TypeScript
root.render(
    <React.StrictMode>
        <App />
    </React.StrictMode>,
);