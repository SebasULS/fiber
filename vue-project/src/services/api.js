import axios from 'axios'

const api = axios.create({
  baseURL: '/api/restful'
})

export const getEnrollmentCertificates = (cui) => {
  return api.get('/enrollment-certificate/', {
    params: { cui }
  })
}

export default api