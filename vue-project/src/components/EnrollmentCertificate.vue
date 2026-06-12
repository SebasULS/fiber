<script setup>
import { ref, onMounted } from 'vue'
import { getEnrollmentCertificates } from '../services/api'

const cui = ref('20250100')
const certificates = ref([])
const student = ref(null)
const loading = ref(false)
const error = ref(null)

const fetchCertificates = async () => {
  if (!cui.value.trim()) return
  loading.value = true
  error.value = null
  certificates.value = []
  student.value = null

  try {
    const response = await getEnrollmentCertificates(cui.value)
    const data = response.data
    certificates.value = data.results
    if (data.results.length > 0) {
      student.value = data.results[0].student
    }
  } catch (err) {
    error.value = err.response?.data?.detail || err.message || 'Error al consultar la API'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchCertificates()
})
</script>

<template>
  <div class="certificate-container">
    <h1>Constancia de Matrícula</h1>

    <div class="search-bar">
      <input
        v-model="cui"
        type="text"
        placeholder="Ingrese CUI del estudiante"
        @keyup.enter="fetchCertificates"
      />
      <button @click="fetchCertificates" :disabled="loading">
        {{ loading ? 'Buscando...' : 'Buscar' }}
      </button>
    </div>

    <div v-if="error" class="error">
      <p>{{ error }}</p>
    </div>

    <div v-if="loading" class="loading">
      <p>Cargando...</p>
    </div>

    <div v-if="student && !loading" class="student-info">
      <h2>Datos del Estudiante</h2>
      <p><strong>Nombre:</strong> {{ student.full_name }}</p>
      <p><strong>CUI:</strong> {{ student.cui }}</p>
      <p><strong>Correo:</strong> {{ student.email }}</p>
    </div>

    <div v-if="certificates.length > 0 && !loading" class="courses">
      <h2>Cursos Matriculados ({{ certificates.length }})</h2>
      <table>
        <thead>
          <tr>
            <th>Código</th>
            <th>Curso</th>
            <th>Créditos</th>
            <th>Grupo</th>
            <th>Laboratorio</th>
            <th>Docente</th>
            <th>Año / Semestre</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="cert in certificates" :key="cert.id">
            <td>{{ cert.workload.course.code }}</td>
            <td>{{ cert.workload.course.name }}</td>
            <td>{{ cert.workload.course.credits }}</td>
            <td>{{ cert.workload.group }}</td>
            <td>{{ cert.workload.laboratory }}</td>
            <td>{{ cert.workload.teacher.full_name }}</td>
            <td>{{ cert.workload.course.year_display }} - {{ cert.workload.course.semester_display }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="certificates.length === 0 && !loading && !error && student === null" class="empty">
      <p>Ingrese un CUI para buscar constancias de matrícula.</p>
    </div>
  </div>
</template>

<style scoped>
.certificate-container {
  max-width: 960px;
  margin: 0 auto;
  padding: 2rem;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}

h1 {
  text-align: center;
  color: #1a237e;
  margin-bottom: 1.5rem;
}

h2 {
  color: #283593;
  border-bottom: 2px solid #3f51b5;
  padding-bottom: 0.5rem;
  margin-bottom: 1rem;
}

.search-bar {
  display: flex;
  gap: 0.5rem;
  justify-content: center;
  margin-bottom: 2rem;
}

.search-bar input {
  padding: 0.6rem 1rem;
  border: 1px solid #90a4ae;
  border-radius: 4px;
  font-size: 1rem;
  width: 260px;
}

.search-bar button {
  padding: 0.6rem 1.5rem;
  background-color: #3f51b5;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  cursor: pointer;
  transition: background-color 0.2s;
}

.search-bar button:hover:not(:disabled) {
  background-color: #303f9f;
}

.search-bar button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.student-info {
  background: #e8eaf6;
  border-radius: 8px;
  padding: 1.2rem 1.5rem;
  margin-bottom: 2rem;
}

.student-info p {
  margin: 0.4rem 0;
  font-size: 1rem;
}

.courses {
  margin-top: 1rem;
}

table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 0.5rem;
}

thead {
  background-color: #3f51b5;
  color: white;
}

th, td {
  padding: 0.7rem 0.8rem;
  text-align: left;
  border-bottom: 1px solid #c5cae9;
}

tbody tr:hover {
  background-color: #e8eaf6;
}

.error {
  background: #ffebee;
  color: #c62828;
  padding: 1rem;
  border-radius: 6px;
  text-align: center;
}

.loading {
  text-align: center;
  padding: 2rem;
  color: #5c6bc0;
  font-size: 1.1rem;
}

.empty {
  text-align: center;
  color: #78909c;
  margin-top: 2rem;
}
</style>