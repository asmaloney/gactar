import { defineConfig } from 'vite'

import compress from 'vite-plugin-compression'
import vue from '@vitejs/plugin-vue2'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [compress({ algorithm: 'brotliCompress' }), vue()],

  server: {
    // This lets us run gactar to serve the endpoints, but run the UI through
    // vite for testing. When running "npm run dev", the frontend will be updated
    // live and the backend will be served by running "gactar -w".
    proxy: {
      '/api': 'http://localhost:8181',
    },
  },
})
