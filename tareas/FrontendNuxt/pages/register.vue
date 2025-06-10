<template>
  <div class="app-wrapper">
    <div class="main-container">
      <div class="header">
        <div class="header-center">
          <h1>Crear Cuenta</h1>
          <p>Únete para organizar tus tareas</p>
        </div>
      </div>

      <div class="form-section">
        <form @submit.prevent="handleRegister" class="auth-form">
          <div class="form-group">
            <label for="name">Nombre completo</label>
            <input
              id="name"
              v-model="form.name"
              type="text"
              required
              placeholder="Tu nombre completo"
            />
          </div>

          <div class="form-group">
            <label for="email">Email</label>
            <input
              id="email"
              v-model="form.email"
              type="email"
              required
              placeholder="tu@email.com"
            />
          </div>

          <div class="form-group">
            <label for="password">Contraseña</label>
            <input
              id="password"
              v-model="form.password"
              type="password"
              required
              placeholder="Mínimo 6 caracteres"
              minlength="6"
            />
          </div>

          <div class="form-group">
            <label for="confirmPassword">Confirmar contraseña</label>
            <input
              id="confirmPassword"
              v-model="form.confirmPassword"
              type="password"
              required
              placeholder="Confirma tu contraseña"
            />
          </div>

          <button
            type="submit"
            :disabled="loading || !passwordsMatch"
            class="btn btn-primary"
          >
            <span v-if="loading">Creando cuenta...</span>
            <span v-else>Crear Cuenta</span>
          </button>
        </form>

        <div v-if="error" class="message error">
          {{ error }}
        </div>

        <div v-if="!passwordsMatch && form.confirmPassword" class="message error">
          Las contraseñas no coinciden
        </div>

        <div v-if="success" class="message success">
          Cuenta creada exitosamente. Redirigiendo a inicio de sesión...
        </div>

        <div class="auth-footer">
          <p>¿Ya tienes cuenta? 
            <NuxtLink to="/login" class="link">Inicia sesión aquí</NuxtLink>
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
const form = ref({
  name: '',
  email: '',
  password: '',
  confirmPassword: ''
})

const loading = ref(false)
const error = ref('')
const success = ref(false)

const passwordsMatch = computed(() => {
  return form.value.password === form.value.confirmPassword
})

// Description: Method to handle user registration
const handleRegister = async () => {
  if (!passwordsMatch.value) {
    error.value = 'Las contraseñas no coinciden'
    return
  }

  loading.value = true
  error.value = ''
  
  try {
    const response = await $fetch('http://localhost:8080/auth/register', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        name: form.value.name,
        email: form.value.email,
        password: form.value.password
      })
    })
    success.value = true
    
    form.value = {
      name: '',
      email: '',
      password: '',
      confirmPassword: ''
    }
    
    setTimeout(() => {
      navigateTo('/login')
    }, 2000)
    
  } catch (err) {
    console.error('Error registering:', err)
    error.value = err.data?.error || 'Error al crear la cuenta'
  } finally {
    loading.value = false
  }
}

// Description: Redirect to home if user is already logged in
onMounted(() => {
  const token = localStorage.getItem('token')
  if (token) {
    navigateTo('/')
  }
})

// Description: Set page title for SEO
useHead({
  title: 'Crear Cuenta - Task Manager'
})
</script>

<style scoped>
.app-wrapper {
  min-height: 100vh;
  padding: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.main-container {
  max-width: 550px;
  width: 100%;
  margin: 0 auto;
  background: white;
  border-radius: 20px;
  padding: 30px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
  border: 2px solid #e9ecef;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}

.header {
  margin-bottom: 30px;
  padding-bottom: 15px;
  border-bottom: 2px solid #e9ecef;
  text-align: center;
}

.header-center h1 {
  margin: 0 0 10px 0;
  color: #333;
  font-size: 2em;
}

.header-center p {
  margin: 0;
  color: #666;
  font-size: 1em;
  font-weight: normal;
}

.form-section {
  background: #f8f9fa;
  padding: 25px;
  border-radius: 10px;
  border: 1px solid #dee2e6;
  box-shadow: 0 3px 10px rgba(0,0,0,0.1);
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 600;
  color: #495057;
  font-size: 0.95em;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}

.form-group input {
  width: 100%;
  padding: 14px 16px;
  border: 2px solid #e9ecef;
  border-radius: 8px;
  font-size: 16px;
  transition: all 0.3s ease;
  box-sizing: border-box;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}

.form-group input:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.btn {
  background-color: #007bff;
  color: white;
  border: none;
  padding: 14px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.95em;
  transition: background-color 0.2s ease;
  width: 100%;
}

.btn:hover {
  background-color: #0056b3;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.message {
  margin: 20px 0;
  padding: 15px;
  border-radius: 8px;
  text-align: center;
  font-weight: 500;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}

.message.error {
  background-color: #f8d7da;
  color: #721c24;
  border: 1px solid #f5c6cb;
}

.message.success {
  background-color: #d4edda;
  color: #155724;
  border: 1px solid #c3e6cb;
}

.auth-footer {
  text-align: center;
  padding-top: 20px;
  border-top: 1px solid #e9ecef;
}

.auth-footer p {
  margin: 0;
  color: #666;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}

.link {
  color: #667eea;
  text-decoration: none;
  font-weight: 600;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}

.link:hover {
  text-decoration: underline;
}
</style>