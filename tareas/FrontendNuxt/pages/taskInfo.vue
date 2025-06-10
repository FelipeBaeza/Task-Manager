<template>
  <div class="app-wrapper">
    <div class="main-container">
      <div class="header">
        <div class="header-left">
          <h1>Detalles de la Tarea</h1>
          <p v-if="task">{{ task.title }}</p>
        </div>
        <div class="header-actions">
          <button @click="goBack" class="btn btn-secondary">
            Volver
          </button>
        </div>
      </div>

      <div v-if="loading" class="loading">
        <p>Cargando tarea...</p>
      </div>

      <div v-else-if="error" class="message error">
        {{ error }}
      </div>

      <div v-else-if="task" class="task-details">
        <div class="task-header">
          <h2>{{ task.title }}</h2>
          <div class="badges">
            <span :class="getStatusClass(task.status)">
              {{ getStatusText(task.status) }}
            </span>
            <span :class="getPriorityClass(task.priority)">
              {{ getPriorityText(task.priority) }}
            </span>
          </div>
        </div>

        <div class="detail-section">
          <h3>Descripción</h3>
          <p class="description">{{ task.description || 'Sin descripción' }}</p>
        </div>

        <div class="detail-section">
          <h3>Información de Fechas</h3>
          <div class="date-info">
            <p><strong>Creada:</strong> {{ formatDate(task.createdAt) }}</p>
            <p v-if="task.updatedAt && task.updatedAt !== task.createdAt">
              <strong>Última actualización:</strong> {{ formatDate(task.updatedAt) }}
            </p>
          </div>
        </div>

        <div class="detail-section">
          <h3>Asignación</h3>
          <div class="assignment-info">
            <p><strong>Asignada a:</strong> {{ assignedUserName || 'Cargando...' }}</p>
            <p><strong>Creada por:</strong> {{ createdByUserName || 'Cargando...' }}</p>
          </div>
        </div>

        <div class="actions">
          <div class="status-buttons">
            <button 
              v-if="task.status === 'pendiente'"
              @click="updateStatus('en_progreso')"
              class="btn btn-warning"
            >
              Iniciar Tarea
            </button>
            
            <button 
              v-if="task.status === 'en_progreso'"
              @click="updateStatus('completada')"
              class="btn btn-success"
            >
              Completar Tarea
            </button>
            
            <button 
              v-if="task.status === 'completada'"
              @click="updateStatus('pendiente')"
              class="btn btn-secondary"
            >
              Reabrir Tarea
            </button>
            
            <button 
              v-if="task.status === 'completada'"
              @click="deleteTask"
              class="btn btn-danger"
            >
              Eliminar Tarea
            </button>
          </div>
        </div>

        <TaskComments 
          :taskId="taskId" 
          :users="users"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import TaskComments from '~/components/TaskComments.vue'

// Description: This page displays the details of a specific task, including its title, description, status, priority, assigned user, and creation date. It allows users to update the task status and delete the task.
definePageMeta({
  middleware: 'auth'
})

const route = useRoute()
const taskId = route.query.id  // Cambiar de route.params.id a route.query.id

const task = ref(null)
const loading = ref(true)
const error = ref('')
const assignedUserName = ref('')
const createdByUserName = ref('')
const users = ref([])

// Description: Load task details on component mount
onMounted(() => {
  loadTask()
})

// Description: Load task details by ID
const loadTask = async () => {
  try {
    loading.value = true
    error.value = ''
    
    const token = localStorage.getItem('token')
    
    const data = await $fetch(`http://localhost:8080/tasks/${taskId}`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    
    task.value = data
    await loadUsers() 
    
  } catch (err) {
    console.error('Error loading task:', err)
    error.value = 'Error al cargar la tarea'
  } finally {
    loading.value = false
  }
}

// Description: Load users for task assignment
const loadUsers = async () => {
  try {
    const token = localStorage.getItem('token')
    const response = await $fetch('http://localhost:8080/users', {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    
    users.value = response.users
    
    if (task.value) {
      const assignedUser = users.value.find(user => user._id === task.value.assignedTo)
      const createdByUser = users.value.find(user => user._id === task.value.createdBy)
      
      assignedUserName.value = assignedUser ? assignedUser.name : 'Usuario no encontrado'
      createdByUserName.value = createdByUser ? createdByUser.name : 'Usuario no encontrado'
    }
    
  } catch (err) {
    console.error('Error loading users:', err)
  }
}

//Description: Format date to a readable string
const formatDate = (dateTime) => {
  if (!dateTime) return 'N/A'
  
  let date
  if (typeof dateTime === 'object' && dateTime.$date) {
    date = new Date(dateTime.$date)
  } else {
    date = new Date(dateTime)
  }
  
  return date.toLocaleDateString('es-ES', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// Description: Get status text and class based on task status
const getStatusText = (status) => {
  const statusMap = {
    'pendiente': 'Pendiente',
    'en_progreso': 'En Progreso',
    'completada': 'Completada'
  }
  return statusMap[status] || status
}

// Description: Get status class based on task status
const getStatusClass = (status) => {
  return `status status-${status}`
}

// Description: Get priority text based on task priority
const getPriorityText = (priority) => {
  const priorityMap = {
    'alta': 'Alta',
    'media': 'Media',
    'baja': 'Baja'
  }
  return priorityMap[priority] || priority
}

// Description: Get priority text and class based on task priority
const getPriorityClass = (priority) => {
  return `priority priority-${priority}`
}

// Description: Update task status
const updateStatus = async (newStatus) => {
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
    
    task.value.status = newStatus
    task.value.updatedAt = new Date().toISOString()
    
  } catch (error) {
    console.error('Error al actualizar tarea:', error)
  }
}

// Description: Delete task
const deleteTask = async () => {
  if (confirm('¿Estás seguro de que quieres eliminar esta tarea? Esta acción no se puede deshacer.')) {
    try {
      const token = localStorage.getItem('token')
      
      await $fetch(`http://localhost:8080/tasks/${taskId}`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${token}`
        }
      })
      
      navigateTo('/')
      
    } catch (error) {
      console.error('Error al eliminar tarea:', error)
    }
  }
}

// Description: Navigate back to the task list
const goBack = () => {
  navigateTo('/')
}

useHead({
  title: 'Detalles de Tarea - Task Manager'
})
</script>

<style scoped>
.app-wrapper {
  min-height: 100vh;
  padding: 20px;
}

.main-container {
  max-width: 800px;
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

.header-left {
  display: flex;
  flex-direction: column;
  gap: 5px;
  flex: 1;
  min-width: 0; 
}

.header-left h1 {
  margin: 0;
  color: #333;
  font-size: 2em;
}

.header-left p {
  margin: 0;
  color: #666;
  font-size: 1em;
  font-weight: normal;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.header-actions {
  display: flex;
  gap: 10px;
  flex-shrink: 0;
}

.btn {
  background-color: #007bff;
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.95em;
  transition: all 0.2s ease;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}

.btn:hover {
  background-color: #0056b3;
  transform: translateY(-1px);
}

.btn-secondary {
  background-color: #6c757d;
  color: white;
}

.btn-secondary:hover {
  background-color: #5a6268;
}

.btn-warning {
  background-color: #ffc107;
  color: #212529;
}

.btn-warning:hover {
  background-color: #e0a800;
}

.btn-success {
  background-color: #28a745;
}

.btn-success:hover {
  background-color: #218838;
}

.btn-danger {
  background-color: #dc3545;
}

.btn-danger:hover {
  background-color: #c82333;
}

.loading {
  text-align: center;
  padding: 40px;
  color: #666;
}

.message {
  padding: 15px;
  border-radius: 6px;
  text-align: center;
  font-weight: 500;
}

.message.error {
  background-color: #f8d7da;
  color: #721c24;
  border: 1px solid #f5c6cb;
}

.task-details {
  background: #f8f9fa;
  padding: 25px;
  border-radius: 10px;
  border: 1px solid #dee2e6;
}

.task-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 25px;
  padding-bottom: 15px;
  border-bottom: 2px solid #dee2e6;
}

.task-header h2 {
  margin: 0;
  color: #333;
  font-size: 1.8em;
  flex: 1;
}

.badges {
  display: flex;
  flex-direction: column;
  gap: 8px;
  align-items: flex-end;
}

.status, .priority {
  padding: 6px 12px;
  border-radius: 6px;
  font-size: 0.9em;
  font-weight: bold;
  text-align: center;
  min-width: 100px;
}

.status-pendiente {
  background: rgba(255, 193, 7, 0.2);
  color: #856404;
}

.status-en_progreso {
  background: rgba(0, 123, 255, 0.2);
  color: #004085;
}

.status-completada {
  background: rgba(40, 167, 69, 0.2);
  color: #155724;
}

.priority-alta {
  background: rgba(220, 53, 69, 0.2);
  color: #721c24;
}

.priority-media {
  background: rgba(255, 193, 7, 0.2);
  color: #856404;
}

.priority-baja {
  background: rgba(108, 117, 125, 0.2);
  color: #495057;
}

.detail-section {
  margin-bottom: 25px;
}

.detail-section h3 {
  margin: 0 0 10px 0;
  color: #495057;
  font-size: 1.2em;
  border-bottom: 1px solid #dee2e6;
  padding-bottom: 5px;
}

.description {
  color: #555;
  line-height: 1.6;
  font-size: 1.1em;
  margin: 0;
  padding: 10px;
  background: white;
  border-radius: 6px;
  border: 1px solid #dee2e6;
}

.date-info p,
.assignment-info p {
  margin: 8px 0;
  color: #555;
  font-size: 1em;
}

.date-info strong,
.assignment-info strong {
  color: #333;
}

.actions {
  margin-top: 30px;
  padding-top: 20px;
  border-top: 2px solid #dee2e6;
}

.status-buttons {
  display: flex;
  gap: 15px;
  flex-wrap: wrap;
  justify-content: center;
}
</style>