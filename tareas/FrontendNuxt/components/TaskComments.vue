<template>
  <div class="comments-section">
    <div v-if="userId">
      <h3>Comentarios</h3>
      <div class="add-comment">
        <textarea 
          v-model="newComment"
          placeholder="Agregar un comentario..."
          rows="3"
          class="comment-input"
        ></textarea>
        <button 
          @click="addComment"
          :disabled="!newComment.trim() || loading"
          class="btn btn-primary"
        >
          {{ loading ? 'Agregando...' : 'Agregar Comentario' }}
        </button>
      </div>

      <div v-if="comments.length === 0" class="no-comments">
        <p>No hay comentarios aún</p>
      </div>
      
      <div v-else class="comments-list">
        <div 
          v-for="comment in comments" 
          :key="comment._id"
          class="comment-item"
        >
          <div class="comment-header">
            <span class="comment-author">{{ getAuthorName(comment.authorId) }}</span>
            <span class="comment-date">{{ formatDate(comment.createdAt) }}</span>
            <button 
              v-if="canDeleteComment(comment)"
              @click="deleteComment(comment._id)"
              class="btn btn-sm btn-danger"
            >
              Eliminar
            </button>
          </div>
          <p class="comment-text">{{ comment.text }}</p>
        </div>
      </div>
    </div>
    <div v-else class="loading">
      Cargando comentarios...
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  taskId: {
    type: String,
    required: true
  },
  users: {
    type: Array,
    default: () => []
  }
})

const { userId } = useAuth()

const comments = ref([])
const newComment = ref('')
const loading = ref(false)

// Load comments
const loadComments = async () => {
  try {
    const token = localStorage.getItem('token')
    const data = await $fetch(`http://localhost:8080/tasks/${props.taskId}/comments`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    comments.value = data
  } catch (error) {
    console.error('Error loading comments:', error)
  }
}

// Add new comment
const addComment = async () => {
  if (!newComment.value.trim()) return
  
  try {
    loading.value = true
    const token = localStorage.getItem('token')
    
    const data = await $fetch(`http://localhost:8080/tasks/${props.taskId}/comments`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`
      },
      body: {
        text: newComment.value.trim()
      }
    })
    
    comments.value.push(data)
    newComment.value = ''
    
  } catch (error) {
    console.error('Error adding comment:', error)
  } finally {
    loading.value = false
  }
}

// Delete comment
const deleteComment = async (commentId) => {
  if (!confirm('¿Estás seguro de que quieres eliminar este comentario?')) return
  
  try {
    const token = localStorage.getItem('token')
    
    await $fetch(`http://localhost:8080/comments/${commentId}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    
    comments.value = comments.value.filter(c => c.id !== commentId)
    
  } catch (error) {
    console.error('Error deleting comment:', error)
  }
}

// Description: Get author name by ID
const getAuthorName = (authorId) => {
  const author = props.users.find(user => user._id === authorId)
  return author ? author.name : 'Usuario desconocido'
}

// Description: Check if the current user can delete a comment
const canDeleteComment = (comment) => {
  const { userId } = useAuth()
  return userId?.value && comment.createdBy === userId.value
}

// Description: Format date to a readable string
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

watch(userId, (newUserId) => {
  if (newUserId) {
    loadComments()
  }
}, { immediate: true })
</script>

<style scoped>
.comments-section {
  margin-top: 30px;
  padding-top: 20px;
  border-top: 2px solid #dee2e6;
}

.comments-section h3 {
  margin: 0 0 20px 0;
  color: #495057;
  font-size: 1.2em;
}

.add-comment {
  margin-bottom: 20px;
}

.comment-input {
  width: 100%;
  padding: 10px;
  border: 1px solid #ced4da;
  border-radius: 6px;
  resize: vertical;
  font-family: inherit;
  margin-bottom: 10px;
}

.comment-input:focus {
  outline: none;
  border-color: #007bff;
  box-shadow: 0 0 0 2px rgba(0, 123, 255, 0.25);
}

.btn {
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.9em;
  transition: all 0.2s ease;
}

.btn-primary {
  background-color: #007bff;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background-color: #0056b3;
}

.btn-primary:disabled {
  background-color: #6c757d;
  cursor: not-allowed;
}

.no-comments {
  text-align: center;
  color: #6c757d;
  font-style: italic;
  padding: 20px;
}

.comments-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.comment-item {
  background: white;
  border: 1px solid #dee2e6;
  border-radius: 6px;
  padding: 15px;
}

.comment-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.comment-author {
  font-weight: 600;
  color: #495057;
}

.comment-date {
  font-size: 0.85em;
  color: #6c757d;
}

.btn-delete {
  background: none;
  border: none;
  color: #dc3545;
  font-size: 1.2em;
  cursor: pointer;
  padding: 0;
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 3px;
}

.btn-delete:hover {
  background-color: #dc3545;
  color: white;
}

.comment-text {
  margin: 0;
  color: #495057;
  line-height: 1.4;
}

.loading {
  text-align: center;
  color: #007bff;
  font-weight: 500;
  padding: 20px;
}
</style>