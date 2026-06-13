import axios from 'axios'

// En local usamos el proxy de Vite, en producción usamos un proxy CORS intermediario
const baseConfigURL = import.meta.env.PROD
  ? 'https://api.allorigins.win/raw?url=https://sisacad-enrollments-backend.vercel.app/restful'
  : '/api/restful'

const api = axios.create({
  baseURL: baseConfigURL
})

export const getEnrollmentCertificates = (cui) => {
  return api.get('/enrollment-certificate/', {
    params: { cui }
  })
}

export default api