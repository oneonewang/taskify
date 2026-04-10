<template>
  <div v-if="visible" class="dialog-overlay" @click.self="close">
    <div class="dialog">
      <h3>创建项目</h3>

      <form @submit.prevent="handleSubmit">
        <div class="form-group">
          <label>项目名称 *</label>
          <input
            v-model="form.name"
            type="text"
            required
            maxlength="100"
            placeholder="输入项目名称"
          />
        </div>

        <div class="form-group">
          <label>描述</label>
          <textarea
            v-model="form.description"
            rows="3"
            maxlength="500"
            placeholder="输入项目描述（可选）"
          ></textarea>
        </div>

        <div v-if="error" class="error-message">{{ error }}</div>

        <div class="dialog-actions">
          <button type="button" @click="close" class="btn btn-secondary">取消</button>
          <button type="submit" class="btn btn-primary" :disabled="loading">
            {{ loading ? '创建中...' : '创建' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { createProject } from '../../api/project'
import type { Project, ApiResponse } from '../../api/project'

const props = defineProps<{
  visible: boolean
}>()

const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void
  (e: 'success', project: Project): void
}>()

const form = ref({
  name: '',
  description: ''
})
const loading = ref(false)
const error = ref<string | null>(null)

watch(() => props.visible, (newVal) => {
  if (newVal) {
    form.value = { name: '', description: '' }
    error.value = null
  }
})

function close() {
  emit('update:visible', false)
}

async function handleSubmit() {
  if (!form.value.name.trim()) {
    error.value = '请输入项目名称'
    return
  }

  loading.value = true
  error.value = null

  try {
    const res = await createProject({
      name: form.value.name.trim(),
      description: form.value.description.trim()
    }) as ApiResponse<Project>

    if (res.success) {
      emit('success', res.data)
    } else {
      error.value = res.error?.message || '创建失败'
    }
  } catch (e: any) {
    error.value = e.message || '创建失败'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0,0,0,0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.dialog {
  background: #fff;
  border-radius: 8px;
  padding: 24px;
  min-width: 400px;
  max-width: 500px;
}

.dialog h3 {
  margin: 0 0 20px;
  font-size: 18px;
  color: #333;
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  margin-bottom: 6px;
  font-size: 14px;
  color: #666;
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  box-sizing: border-box;
}

.form-group input:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #1976d2;
}

.form-group textarea {
  resize: vertical;
  min-height: 80px;
}

.error-message {
  padding: 8px 12px;
  background: #fef2f2;
  color: #dc3545;
  border-radius: 4px;
  font-size: 14px;
  margin-bottom: 16px;
}

.dialog-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 20px;
}

.btn {
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

.btn-primary {
  background: #1976d2;
  color: #fff;
}

.btn-primary:hover:not(:disabled) {
  background: #1565c0;
}

.btn-primary:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.btn-secondary {
  background: #f5f5f5;
  color: #666;
}

.btn-secondary:hover {
  background: #e0e0e0;
}
</style>
