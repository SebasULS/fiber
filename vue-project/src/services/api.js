import axios from 'axios'

const api = axios.create({
  baseURL: 'https://sisacad-enrollments-backend.vercel.app/restful'
})

export const getEnrollmentCertificates = (cui) => {
  return api.get('/enrollment-certificate/', {
    params: { cui }
  })
}

export default api