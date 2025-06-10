<template>
  <div class="task-card">
    <div class="task-header">
      <h3 class="task-title">{{ task.title }}</h3>
      <div class="badges">
        <span :class="getStatusClass(task.status)">
          {{ getStatusText(task.status) }}
        </span>
        <span :class="getPriorityClass(task.priority)">
          {{ getPriorityText(task.priority) }}
        </span>
      </div>
    </div>
    
    <p>{{ task.description || 'Sin descripción' }}</p>
    
    <div class="dates">
      <p class="date">Creada: {{ formatDate(task.createdAt) }}</p>
      <p v-if="task.updatedAt && task.updatedAt !== task.createdAt" class="date">
        Actualizada: {{ formatDate(task.updatedAt) }}
      </p>
    </div>
    
    <div class="actions">
      <div class="status-buttons">
        <button 
          @click="viewTask"
          class="btn btn-info"
        >
          Ver Detalles
        </button>
        
        <button 
          v-if="task.status === 'pendiente'"
          @click="updateStatus('en_progreso')"
          class="btn btn-warning"
        >
          Iniciar
        </button>
        
        <button 
          v-if="task.status === 'en_progreso'"
          @click="updateStatus('completada')"
          class="btn btn-success"
        >
          Completar
        </button>
        
        <button 
          v-if="task.status === 'completada'"
          @click="updateStatus('pendiente')"
          class="btn btn-secondary"
        >
          Reabrir
        </button>
        
        <button 
          v-if="task.status === 'completada'"
          @click="deleteTask"
          class="btn btn-danger"
        >
          Eliminar
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  task: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['updateStatus', 'deleteTask'])

// Description: Function to format date to a readable string
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
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

//Description: Function to navigate to task details
const getStatusText = (status) => {
  const statusMap = {
    'pendiente': 'Pendiente',
    'en_progreso': 'En Progreso',
    'completada': 'Completada'
  }
  return statusMap[status] || status
}

// Description: Function to get status class based on task status
const getStatusClass = (status) => {
  return `status status-${status}`
}

// Description: Function to get priority text based on task priority
const getPriorityText = (priority) => {
  const priorityMap = {
    'alta': 'Alta',
    'media': 'Media',
    'baja': 'Baja'
  }
  return priorityMap[priority] || priority
}

// Description: Function to get priority class based on task priority
const getPriorityClass = (priority) => {
  return `priority priority-${priority}`
}

// Description: Function to update task status
const updateStatus = (newStatus) => {
  emit('updateStatus', props.task.id, newStatus)
}

// Description: Function to delete task
const deleteTask = () => {
  if (confirm('¿Estás seguro de que quieres eliminar esta tarea? Esta acción no se puede deshacer.')) {
    emit('deleteTask', props.task.id)
  }
}

// Description: Function to navigate to task details page
const viewTask = () => {
  navigateTo(`/taskInfo?id=${props.task.id}`)
}
</script>

<style scoped>
.task-card {
  background: white;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 15px;
  text-align: left;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.task-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.task-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 15px;
}

.task-title {
  margin: 0;
  color: #333;
  font-size: 1.2em;
  flex: 1;
  font-weight: 600;
}

.badges {
  display: flex;
  flex-direction: column;
  gap: 5px;
  align-items: flex-end;
}

.status, .priority {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 0.8em;
  font-weight: bold;
  text-align: center;
  min-width: 80px;
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

.task-card p {
  margin: 8px 0;
  color: #555;
  line-height: 1.4;
}

.dates {
  margin: 10px 0;
}

.date {
  font-size: 0.85em;
  color: #888;
  font-style: italic;
  margin: 2px 0;
}

.actions {
  margin-top: 15px;
  padding-top: 15px;
  border-top: 1px solid #f0f0f0;
}

.status-buttons {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.btn {
  border: none;
  padding: 8px 16px;
  border-radius: 5px;
  cursor: pointer;
  font-size: 0.9em;
  transition: background-color 0.2s ease;
  font-weight: 500;
}

.btn-info {
  background-color: #17a2b8;
  color: white;
}

.btn-info:hover {
  background-color: #138496;
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
  color: white;
}

.btn-success:hover {
  background-color: #218838;
}

.btn-secondary {
  background-color: #6c757d;
  color: white;
}

.btn-secondary:hover {
  background-color: #5a6268;
}

.btn-danger {
  background-color: #dc3545;
  color: white;
}

.btn-danger:hover {
  background-color: #c82333;
}
</style>