<template>
  <div v-if="visible" class="dialog-overlay" @click.self="close">
    <div class="dialog">
      <h3>批量导入用户</h3>

      <div class="dialog-content">
        <p class="description">
          请上传 Excel 文件（.xlsx 格式）<br />
          文件应包含三列：邮箱、显示名称、角色
        </p>

        <div class="upload-area" @click="triggerFileInput">
          <input
            ref="fileInput"
            type="file"
            accept=".xlsx"
            @change="handleFileSelect"
            style="display: none"
          />
          <div v-if="!selectedFile" class="upload-placeholder">
            <span class="upload-icon">📁</span>
            <span>点击选择文件或拖拽到此处</span>
          </div>
          <div v-else class="file-selected">
            <span class="file-icon">📄</span>
            <span>{{ selectedFile.name }}</span>
            <button @click.stop="clearFile" class="btn-clear">×</button>
          </div>
        </div>

        <div v-if="importResult" class="import-result">
          <div class="result-summary">
            <span class="result-item success">成功: {{ importResult.success }}</span>
            <span class="result-item failed">失败: {{ importResult.failed }}</span>
            <span class="result-item total">总计: {{ importResult.total }}</span>
          </div>

          <div v-if="importResult.errors && importResult.errors.length > 0" class="error-list">
            <h4>错误列表：</h4>
            <ul>
              <li v-for="(err, idx) in importResult.errors.slice(0, 10)" :key="idx">{{ err }}</li>
            </ul>
            <p v-if="importResult.errors.length > 10" class="more-errors">
              还有 {{ importResult.errors.length - 10 }} 个错误...
            </p>
          </div>
        </div>

        <div v-if="importing" class="importing">
          <span class="spinner"></span>
          正在导入...
        </div>

        <div v-if="importError" class="import-error">
          {{ importError }}
        </div>
      </div>

      <div class="dialog-actions">
        <button @click="downloadTemplate" class="btn btn-secondary">下载模板</button>
        <button @click="close" class="btn btn-secondary">关闭</button>
        <button
          @click="startImport"
          class="btn btn-primary"
          :disabled="!selectedFile || importing"
        >
          开始导入
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { importUsers, getImportTemplateUrl, type ImportResult } from '../../api/user'

defineProps<{
  visible: boolean
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'success'): void
}>()

const fileInput = ref<HTMLInputElement | null>(null)
const selectedFile = ref<File | null>(null)
const importing = ref(false)
const importResult = ref<ImportResult | null>(null)
const importError = ref<string | null>(null)

function triggerFileInput() {
  fileInput.value?.click()
}

function handleFileSelect(event: Event) {
  const target = event.target as HTMLInputElement
  if (target.files && target.files[0]) {
    selectedFile.value = target.files[0]
    importResult.value = null
    importError.value = null
  }
}

function clearFile() {
  selectedFile.value = null
  if (fileInput.value) {
    fileInput.value.value = ''
  }
}

function downloadTemplate() {
  window.location.href = getImportTemplateUrl()
}

async function startImport() {
  if (!selectedFile.value) return

  importing.value = true
  importError.value = null
  importResult.value = null

  try {
    const res = await importUsers(selectedFile.value)
    if (res.code === 0) {
      importResult.value = res.data
      emit('success')
    } else {
      importError.value = res.message
    }
  } catch (e: any) {
    importError.value = e.response?.data?.message || '导入失败'
  } finally {
    importing.value = false
  }
}

function close() {
  selectedFile.value = null
  importResult.value = null
  importError.value = null
  emit('close')
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
  min-width: 480px;
  max-width: 560px;
}

.dialog h3 {
  margin: 0 0 16px;
  font-size: 18px;
}

.dialog-content {
  margin: 16px 0;
}

.description {
  color: #666;
  font-size: 14px;
  line-height: 1.5;
  margin-bottom: 16px;
}

.upload-area {
  border: 2px dashed #ddd;
  border-radius: 8px;
  padding: 32px;
  text-align: center;
  cursor: pointer;
  transition: border-color 0.2s;
}

.upload-area:hover {
  border-color: #1976d2;
}

.upload-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  color: #666;
}

.upload-icon {
  font-size: 32px;
}

.file-selected {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  color: #333;
}

.file-icon {
  font-size: 24px;
}

.btn-clear {
  background: none;
  border: none;
  font-size: 20px;
  cursor: pointer;
  color: #999;
  padding: 0 4px;
}

.btn-clear:hover {
  color: #dc3545;
}

.import-result {
  margin-top: 16px;
  padding: 12px;
  background: #f8f9fa;
  border-radius: 4px;
}

.result-summary {
  display: flex;
  gap: 16px;
  justify-content: center;
}

.result-item {
  font-weight: 600;
}

.result-item.success {
  color: #28a745;
}

.result-item.failed {
  color: #dc3545;
}

.result-item.total {
  color: #333;
}

.error-list {
  margin-top: 12px;
  text-align: left;
}

.error-list h4 {
  margin: 8px 0;
  font-size: 14px;
}

.error-list ul {
  margin: 0;
  padding-left: 20px;
  font-size: 12px;
  color: #dc3545;
}

.error-list li {
  margin: 4px 0;
}

.more-errors {
  font-size: 12px;
  color: #666;
  margin-top: 8px;
}

.importing {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  margin-top: 16px;
  color: #1976d2;
}

.spinner {
  width: 16px;
  height: 16px;
  border: 2px solid #e3f2fd;
  border-top-color: #1976d2;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.import-error {
  margin-top: 16px;
  padding: 12px;
  background: #ffebee;
  color: #c62828;
  border-radius: 4px;
  font-size: 14px;
}

.dialog-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 16px;
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

.btn-primary:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.btn-secondary {
  background: #6c757d;
  color: #fff;
}
</style>
