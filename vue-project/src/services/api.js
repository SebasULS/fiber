import axios from 'axios'

const api = axios.create({
  baseURL: '/api' // Así limpio, sin el 'restful' aquí
})

export const getEnrollmentCertificates = (cui) => {
  // Ponemos la ruta completa que necesita la API externa desde aquí
  return api.get('/restful/enrollment-certificate/', {
    params: { cui }
  })
}

export default api