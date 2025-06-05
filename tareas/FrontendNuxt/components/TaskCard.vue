<template>
  <div class="task-card">
    <h3>{{ task.title }}</h3>
    <p>{{ task.description || 'Sin descripción' }}</p>
    <p class="date">Creada: {{ formatDate(task.created_at) }}</p>
    
    <div class="footer">
      <span :class="task.completed ? 'completed' : 'pending'">
        {{ task.completed ? 'Completada' : 'Pendiente' }}
      </span>
      
      <button 
        v-if="!task.completed"
        @click="$emit('complete', task.id)"
        class="btn"
      >
        Completar
      </button>
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

defineEmits(['complete'])

const formatDate = (iso) => {
  return new Date(iso).toLocaleDateString('es-ES')
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

.task-card h3 {
  margin: 0 0 10px 0;
  color: #333;
  font-size: 1.2em;
}

.task-card p {
  margin: 8px 0;
  color: #555;
  line-height: 1.4;
}

.date {
  font-size: 0.85em;
  color: #888;
  font-style: italic;
}

.footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 15px;
  padding-top: 15px;
  border-top: 1px solid #f0f0f0;
}

.completed {
  color: #28a745;
  font-weight: bold;
  background: rgba(40, 167, 69, 0.1);
  padding: 4px 8px;
  border-radius: 4px;
}

.pending {
  color: #ffc107;
  font-weight: bold;
  background: rgba(255, 193, 7, 0.1);
  padding: 4px 8px;
  border-radius: 4px;
}

.btn {
  background-color: #007bff;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 5px;
  cursor: pointer;
  font-size: 0.9em;
  transition: background-color 0.2s ease;
}

.btn:hover {
  background-color: #0056b3;
}
</style>