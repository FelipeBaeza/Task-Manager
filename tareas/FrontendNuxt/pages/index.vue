<template>
  <div class="app-wrapper">
    <div class="main-container">
      <div class="header">
        <h1>Lista de Tareas</h1>
        <button @click="createTasks" class="btn">
          Nueva Tarea
        </button>
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
            @complete="completeTask"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import axios from 'axios'
import TaskCard from '~/components/TaskCard.vue'

const tasks = ref([])
const loading = ref(true)

// Description: Method to get all tasks
const loadTask = async () => {
  try {
    const response = await axios.get('http://localhost:8080/tasks')
    tasks.value = response.data
  } catch (error) {
    console.error('Error al cargar tareas:', error)
  } finally {
    loading.value = false
  }
}

// Description: Method to navigate to the create task page
const createTasks = () => {
  navigateTo('/createTask')
}

// Description: Method to complete a task
const completeTask = async (id) => {
  try {
    await axios.put(`http://localhost:8080/tasks/${id}/complete`)
    const tarea = tasks.value.find((t) => t.id === id)
    if (tarea) tarea.completed = true
  } catch (error) {
    console.error('Error al completar tarea:', error)
  }
}

// Description: method that automatically executes tasks
onMounted(loadTask)

// Descrition: Reload tasks
onActivated(() => {
  loadTask()
})
</script>

<style scoped>
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

.btn {
  background-color: #007bff;
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.95em;
  transition: all 0.2s ease;
}

.btn:hover {
  background-color: #0056b3;
  transform: translateY(-1px);
}

.btn-green {
  background-color: #28a745;
}

.btn-green:hover {
  background-color: #1e7e34;
}

.form {
  background: #f8f9fa;
  padding: 25px;
  border-radius: 10px;
  margin-bottom: 30px;
  border: 1px solid #dee2e6;
  box-shadow: 0 3px 10px rgba(0,0,0,0.1);
}

.form div {
  margin-bottom: 18px;
}

.form label {
  display: block;
  margin-bottom: 6px;
  font-weight: 600;
  color: #495057;
  font-size: 0.95em;
}

.form input,
.form textarea {
  width: 100%;
  padding: 12px;
  border: 1px solid #ced4da;
  border-radius: 6px;
  font-size: 14px;
  transition: border-color 0.15s ease-in-out, box-shadow 0.15s ease-in-out;
  font-family: inherit;
}

.form input:focus,
.form textarea:focus {
  outline: none;
  border-color: #80bdff;
  box-shadow: 0 0 0 0.2rem rgba(0,123,255,0.25);
}

.content {
  text-align: center;
}

.content p {
  color: #666;
  font-size: 1.1em;
}
</style>
