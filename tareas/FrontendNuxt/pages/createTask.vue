<template>
  <div class="app-wrapper">
    <div class="main-container">
      <div class="header">
        <h1>Crear Nueva Tarea</h1>
        <button @click="goBack" class="btn">
          Volver
        </button>
      </div>

      <div class="form">
        <form @submit.prevent="createTask">
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

          <button
            type="submit"
            :disabled="loading || !task.title.trim()"
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
const task = ref({
  title: '',
  description: ''
})

const loading = ref(false)
const showSuccess = ref(false)
const error = ref('')


// Description: Method to create a new task and navegate back to the main page
const createTask = async () => {
  loading.value = true
  error.value = ''
  
  try {
    await $fetch('http://localhost:8080/tasks', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(task.value)
    })
    
    showSuccess.value = true
    
    task.value = {
      title: '',
      description: ''
    }
    
    setTimeout(() => {
      navigateTo('/')
    }, 2000)
    
  } catch (err) {
    console.error('Error creating task:', err)
    error.value = err.data?.error || 'Error al crear la tarea'
  } finally {
    loading.value = false
  }
}

// Description: Method to navigate back to the main page
const goBack = () => {
  navigateTo('/')
}

// Description: Set the page title for SEO
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
  transition: background-color 0.2s ease;
}

.btn:hover {
  background-color: #0056b3;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-green {
  background-color: #28a745;
  transition: all 0.2s ease;
}

.btn-green:hover:not(:disabled) {
  background-color: #1e7e34;
  transform: translateY(-1px);
}

.form {
  background: #f8f9fa;
  padding: 25px;
  border-radius: 10px;
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

.required {
  color: #dc3545;
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

.form textarea {
  resize: vertical;
  min-height: 100px;
}

.message {
  margin-top: 20px;
  padding: 15px;
  border-radius: 6px;
  text-align: center;
  font-weight: 500;
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