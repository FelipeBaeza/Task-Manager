<template>
  <div class="app-wrapper">
    <div class="main-container">
      <div class="header">
        <div class="header-left">
          <h1>Crear Nueva Tarea</h1>
          <p>Completa los campos para crear una nueva tarea</p>
        </div>
        <div class="header-actions">
          <button @click="goBack" class="btn btn-secondary">
            Volver
          </button>
        </div>
      </div>

      <div class="form-section">
        <form @submit.prevent="createTask" class="task-form">
          <div>
            <label for="title">
              Título <span class="required">*</span>
            </label>
            <input
              id="title"
              v-model="task.title"
              type="text"
              required
              placeholder="Ingresa el título de la tarea"
            />
          </div>

          <div>
            <label for="description">
              Descripción
            </label>
            <textarea
              id="description"
              v-model="task.description"
              rows="4"
              placeholder="Describe tu tarea (opcional)"
            ></textarea>
          </div>

          <div>
            <label for="priority">
              Prioridad
            </label>
            <select id="priority" v-model="task.priority">
              <option value="baja">Baja</option>
              <option value="media">Media</option>
              <option value="alta">Alta</option>
            </select>
          </div>

          <div>
            <label for="assignedTo">
              Asignar a <span class="required">*</span>
            </label>
            <div class="user-selector">
              <input
                id="assignedTo"
                v-model="searchTerm"
                type="text"
                placeholder="Buscar usuario..."
                @input="filterUsers"
                @focus="showDropdown = true"
                @blur="hideDropdown"
                class="search-input"
              />
              <div v-if="showDropdown && filteredUsers.length > 0" class="dropdown">
                <div
                  v-for="user in filteredUsers"
                  :key="user._id"
                  @mousedown="selectUser(user)"
                  class="dropdown-item"
                >
                  {{ user.name }}
                </div>
              </div>
              <div v-if="selectedUser" class="selected-user">
                <span>{{ selectedUser.name }}</span>
                <button type="button" @click="clearSelection" class="clear-btn">×</button>
              </div>
            </div>
          </div>

          <button
            type="submit"
            :disabled="loading || !task.title.trim() || !selectedUser"
            class="btn btn-green"
          >
            <span v-if="loading">Creando...</span>
            <span v-else>Crear Tarea</span>
          </button>
        </form>

        <div v-if="showSuccess" class="message success">
          ¡Tarea creada exitosamente!
        </div>

        <div v-if="error" class="message error">
          {{ error }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>

definePageMeta({
  middleware: 'auth'
})

const { logout } = useAuth()

const task = ref({
  title: '',
  description: '',
  priority: 'media'
})

const users = ref([])
const filteredUsers = ref([])
const selectedUser = ref(null)
const searchTerm = ref('')
const showDropdown = ref(false)
const loading = ref(false)
const showSuccess = ref(false)
const error = ref('')

// Description: method to load all users from the backend
const loadUsers = async () => {
  try {
    const token = localStorage.getItem('token')
    const response = await $fetch('http://localhost:8080/users', {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    
    users.value = response.users.map(user => ({
      _id: user._id,
      name: user.name
    }))
    filteredUsers.value = users.value
    
  } catch (err) {
    console.error('Error loading users:', err)
    error.value = 'Error al cargar usuarios'
  }
}

// Description: filter users based on search term
const filterUsers = () => {
  if (!searchTerm.value.trim()) {
    filteredUsers.value = users.value
  } else {
    filteredUsers.value = users.value.filter(user =>
      user.name.toLowerCase().includes(searchTerm.value.toLowerCase())
    )
  }
  showDropdown.value = true
}

// Description: select a user from the dropdown
const selectUser = (user) => {
  selectedUser.value = user
  searchTerm.value = user.name
  showDropdown.value = false
}

// Description: clear the selected user and search term
const clearSelection = () => {
  selectedUser.value = null
  searchTerm.value = ''
  showDropdown.value = false
}

// Description: hide the dropdown after a short delay
const hideDropdown = () => {
  setTimeout(() => {
    showDropdown.value = false
  }, 200)
}

// Description: method to create a new task
const createTask = async () => {
  loading.value = true
  error.value = ''
  
  try {
    const token = localStorage.getItem('token')
    
    const taskData = {
      title: task.value.title,
      description: task.value.description,
      priority: task.value.priority,
      assignedTo: selectedUser.value._id
    }
    
    await $fetch('http://localhost:8080/tasks', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`
      },
      body: taskData  
    })
    
    showSuccess.value = true
    
    task.value = {
      title: '',
      description: '',
      priority: 'media'
    }
    clearSelection()
    
    setTimeout(() => {
      navigateTo('/')
    }, 2000)
    
  } catch (err) {
    console.error('Error creating task:', err)
    error.value = err.data?.error || 'Error al crear la tarea'
    if (error.status === 401) {
      logout()
    }
  } finally {
    loading.value = false
  }
}

const goBack = () => {
  navigateTo('/')
}

onMounted(() => {
  loadUsers()
})

useHead({
  title: 'Crear Tarea - Task Manager'
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

.header-left {
  display: flex;
  flex-direction: column;
  gap: 5px;
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
}

.header-actions {
  display: flex;
  gap: 10px;
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
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}

.btn:hover {
  background-color: #0056b3;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-secondary {
  background-color: #6c757d;
  color: white;
}

.btn-secondary:hover {
  background-color: #5a6268;
}

.btn-green {
  background-color: #28a745;
  transition: background-color 0.2s ease;
}

.btn-green:hover:not(:disabled) {
  background-color: #1e7e34;
}

.form-section {
  background: #f8f9fa;
  padding: 25px;
  border-radius: 10px;
  border: 1px solid #dee2e6;
  box-shadow: 0 3px 10px rgba(0,0,0,0.1);
}

.task-form div {
  margin-bottom: 18px;
}

.task-form label {
  display: block;
  margin-bottom: 6px;
  font-weight: 600;
  color: #495057;
  font-size: 0.95em;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}

.required {
  color: #dc3545;
}

.task-form input,
.task-form textarea,
.task-form select {
  width: 100%;
  padding: 12px;
  border: 1px solid #ced4da;
  border-radius: 6px;
  font-size: 14px;
  transition: border-color 0.15s ease-in-out, box-shadow 0.15s ease-in-out;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  box-sizing: border-box;
}

.task-form input:focus,
.task-form textarea:focus,
.task-form select:focus {
  outline: none;
  border-color: #80bdff;
  box-shadow: 0 0 0 0.2rem rgba(0,123,255,0.25);
}

.task-form textarea {
  resize: vertical;
  min-height: 100px;
}

.user-selector {
  position: relative;
}

.search-input {
  width: 100%;
}

.dropdown {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: white;
  border: 1px solid #ced4da;
  border-top: none;
  border-radius: 0 0 6px 6px;
  max-height: 200px;
  overflow-y: auto;
  z-index: 1000;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.dropdown-item {
  padding: 12px;
  cursor: pointer;
  border-bottom: 1px solid #f8f9fa;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}

.dropdown-item:hover {
  background-color: #f8f9fa;
}

.dropdown-item:last-child {
  border-bottom: none;
}

.selected-user {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 8px;
  padding: 8px 12px;
  background-color: #e3f2fd;
  border: 1px solid #90caf9;
  border-radius: 6px;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}

.clear-btn {
  background: none;
  border: none;
  font-size: 18px;
  color: #666;
  cursor: pointer;
  padding: 0;
  margin-left: 8px;
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.clear-btn:hover {
  color: #333;
}

.message {
  margin-top: 20px;
  padding: 15px;
  border-radius: 6px;
  text-align: center;
  font-weight: 500;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}

.message.success {
  background-color: #d4edda;
  color: #155724;
  border: 1px solid #c3e6cb;
}

.message.error {
  background-color: #f8d7da;
  color: #721c24;
  border: 1px solid #f5c6cb;
}
</style>