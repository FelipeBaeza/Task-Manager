<template>
  <div class="app-wrapper">
    <div class="main-container">
      <div class="header">
        <div class="header-left">
          <h1>Lista de Tareas</h1>
          <p v-if="userId" class="welcome">¡Hola!</p>
        </div>
        <div class="header-actions">
          <button @click="createTasks" class="btn">
            Nueva Tarea
          </button>
          <button @click="handleLogout" class="btn btn-secondary">
            Cerrar Sesión
          </button>
        </div>
      </div>

      <!-- Filtros -->
      <div class="filters">
        <div class="filter-group">
          <label>Estado:</label>
          <select v-model="statusFilter" @change="loadTasks">
            <option value="">Todos</option>
            <option value="pendiente">Pendiente</option>
            <option value="en_progreso">En Progreso</option>
            <option value="completada">Completada</option>
          </select>
        </div>
        
        <div class="filter-group">
          <label>Prioridad:</label>
          <select v-model="priorityFilter" @change="loadTasks">
            <option value="">Todas</option>
            <option value="alta">Alta</option>
            <option value="media">Media</option>
            <option value="baja">Baja</option>
          </select>
        </div>
      </div>

      <div class="content">
        <div v-if="loading">
          <p>Cargando tareas...</p>
        </div>
        
        <div v-else-if="tasks.length === 0">
          <p>No hay tareas disponibles</p>
        </div>
        
        <div v-else>
          <TaskCard 
            v-for="task in tasks"
            :key="task.id"
            :task="task"
            @updateStatus="updateTaskStatus"
            @deleteTask="deleteTask"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import TaskCard from '~/components/TaskCard.vue'


definePageMeta({
  middleware: 'auth'
})

const { userId, logout } = useAuth()

const tasks = ref([])
const loading = ref(true)
const statusFilter = ref('')
const priorityFilter = ref('')

// Method to get all tasks with filters
const loadTasks = async () => {
  try {
    loading.value = true
    
    const token = localStorage.getItem('token')
    
    // Build query parameters
    const params = new URLSearchParams()
    if (statusFilter.value) params.append('status', statusFilter.value)
    if (priorityFilter.value) params.append('priority', priorityFilter.value)
    
    const queryString = params.toString()
    const url = queryString ? `http://localhost:8080/tasks?${queryString}` : 'http://localhost:8080/tasks'
    
    const data = await $fetch(url, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    tasks.value = data
  } catch (error) {
    console.error('Error al cargar tareas:', error)
    if (error.status === 401) {
      logout()
    }
  } finally {
    loading.value = false
  }
}

// Description: Method to navigate to the create task page
const createTasks = () => {
  navigateTo('/createTask')
}

// Method to update task status
const updateTaskStatus = async (taskId, newStatus) => {
  try {
    const token = localStorage.getItem('token')
    
    await $fetch(`http://localhost:8080/tasks/${taskId}`, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${token}`
      },
      body: {
        status: newStatus
      }
    })
    
    // Update local task
    const task = tasks.value.find((t) => t.id === taskId)
    if (task) {
      task.status = newStatus
      task.updatedAt = new Date().toISOString()
    }
  } catch (error) {
    console.error('Error al actualizar tarea:', error)
    if (error.status === 401) {
      logout()
    }
  }
}

// Description: Method to delete task
const deleteTask = async (taskId) => {
  try {
    const token = localStorage.getItem('token')
    
    await $fetch(`http://localhost:8080/tasks/${taskId}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    
    // Remove task from local array
    tasks.value = tasks.value.filter(task => task.id !== taskId)
    
  } catch (error) {
    console.error('Error al eliminar tarea:', error)
    if (error.status === 401) {
      logout()
    }
  }
}

const handleLogout = () => {
  logout()
}

// Load tasks on component mount
onMounted(() => {
  loadTasks()
})

// Reload tasks when navigating back
onActivated(() => {
  loadTasks()
})
</script>

<style scoped>
.header-left {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.welcome {
  margin: 0;
  color: #666;
  font-size: 1em;
  font-weight: normal;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.btn-secondary {
  background-color: #6c757d;
  color: white;
}

.btn-secondary:hover {
  background-color: #5a6268;
}

.app-wrapper {
  min-height: 100vh;
  padding: 20px;
}

.main-container {
  max-width: 700px;
  margin: 0 auto;
  background: white;
  border-radius: 20px;
  padding: 30px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
  border: 2px solid #e9ecef;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
  padding-bottom: 15px;
  border-bottom: 2px solid #e9ecef;
}

.header h1 {
  margin: 0;
  color: #333;
  font-size: 2em;
}

.filters {
  display: flex;
  gap: 20px;
  margin-bottom: 20px;
  padding: 15px;
  background: #f8f9fa;
  border-radius: 8px;
}

.filter-group {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.filter-group label {
  font-weight: 600;
  color: #495057;
  font-size: 0.9em;
}

.filter-group select {
  padding: 8px 12px;
  border: 1px solid #ced4da;
  border-radius: 4px;
  font-size: 14px;
}

.btn {
  background-color: #007bff;
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.95em;
  transition: background-color 0.2s ease;
}

.btn:hover {
  background-color: #0056b3;
}

.content {
  text-align: center;
}

.content p {
  color: #666;
  font-size: 1.1em;
}
</style>
