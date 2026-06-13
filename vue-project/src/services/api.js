import axios from 'axios'

const api = axios.create({
  // En local usamos tu proxy de Vite.
  // En producción, usamos un proxy que no se rompe con la estructura de Axios.
  baseURL: import.meta.env.PROD
    ? 'https://cors-anywhere.herokuapp.com/https://sisacad-enrollments-backend.vercel.app/restful'
    : '/api/restful'
})

export const getEnrollmentCertificates = (cui) => {
  return api.get('/enrollment-certificate/', {
    params: { cui }
  })
}

export default api