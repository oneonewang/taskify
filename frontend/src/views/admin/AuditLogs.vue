<template>
  <div class="admin-page">
    <div class="page-header animate-in">
      <div class="header-content">
        <div class="page-title-section">
          <h1>审计日志</h1>
          <p>查看系统操作记录和用户活动</p>
        </div>
      </div>

      <!-- Filter Form -->
      <div class="filter-form">
        <div class="filter-row">
          <div class="filter-item">
            <label>用户ID</label>
            <input
              v-model="filterForm.user_id"
              type="number"
              placeholder="输入用户ID"
              class="filter-input"
              @keyup.enter="handleSearch"
            />
          </div>
          <div class="filter-item">
            <label>事件类型</label>
            <select v-model="filterForm.event_type" class="filter-select" @change="handleSearch">
              <option value="">全部类型</option>
              <option v-for="type in EVENT_TYPES" :key="type.value" :value="type.value">
                {{ type.label }}
              </option>
            </select>
          </div>
          <div class="filter-item date-range">
            <label>时间范围</label>
            <el-date-picker
              v-model="dateRange"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              value-format="YYYY-MM-DD"
              @change="handleDateChange"
              class="date-picker"
            />
          </div>
          <div class="filter-actions">
            <button @click="handleSearch" class="btn btn-primary btn-sm">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="11" cy="11" r="8"/>
                <line x1="21" y1="21" x2="16.65" y2="16.65"/>
              </svg>
              搜索
            </button>
            <button @click="handleReset" class="btn btn-ghost btn-sm">重置</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Logs Table -->
    <div v-if="loading" class="loading-state animate-in">
      <div class="loading-spinner"></div>
    </div>
    <div v-else-if="error" class="error-state animate-in">
      <div class="error-icon">
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="10"/>
          <line x1="12" y1="8" x2="12" y2="12"/>
          <line x1="12" y1="16" x2="12.01" y2="16"/>
        </svg>
      </div>
      <span>{{ error }}</span>
    </div>
    <template v-else>
      <div class="logs-table-container animate-in">
        <table class="logs-table">
          <thead>
            <tr>
              <th>ID</th>
              <th>用户</th>
              <th>事件类型</th>
              <th>详情</th>
              <th>IP地址</th>
              <th>时间</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="log in logs" :key="log.id" class="log-row">
              <td class="id-cell">{{ log.id }}</td>
              <td class="user-cell">
                <div class="user-info">
                  <span class="user-email">{{ log.user_email || `用户ID: ${log.user_id}` }}</span>
                </div>
              </td>
              <td class="event-cell">
                <span :class="['event-tag', getEventTypeClass(log.event_type)]">
                  {{ getEventTypeLabel(log.event_type) }}
                </span>
              </td>
              <td class="details-cell">
                <span class="details-text" :title="log.details">{{ log.details || '-' }}</span>
              </td>
              <td class="ip-cell">{{ log.ip_address || '-' }}</td>
              <td class="time-cell">{{ formatDateTime(log.created_at) }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div class="pagination-wrapper">
        <div class="pagination-info">
          共 {{ total }} 条记录
        </div>
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[20, 50, 100]"
          layout="sizes, prev, pager, next"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { getAuditLogs, EVENT_TYPES, type AuditLog, type AuditLogQuery } from '../../api/audit'

const logs = ref<AuditLog[]>([])
const loading = ref(false)
const error = ref<string | null>(null)
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const dateRange = ref<[string, string] | null>(null)

const filterForm = reactive<{
  user_id?: number
  event_type?: string
}>({})

onMounted(async () => {
  await loadLogs()
})

async function loadLogs() {
  loading.value = true
  error.value = null

  const query: AuditLogQuery = {
    page: currentPage.value,
    page_size: pageSize.value
  }

  if (filterForm.user_id) {
    query.user_id = filterForm.user_id
  }
  if (filterForm.event_type) {
    query.event_type = filterForm.event_type
  }
  if (dateRange.value) {
    query.from = dateRange.value[0]
    query.to = dateRange.value[1]
  }

  try {
    const res = await getAuditLogs(query)
    if (res.code === 0) {
      logs.value = res.data.logs
      total.value = res.data.total
    } else {
      error.value = res.message
    }
  } catch (e: any) {
    error.value = e.response?.data?.message || e.message || '加载失败'
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  currentPage.value = 1
  loadLogs()
}

function handleReset() {
  filterForm.user_id = undefined
  filterForm.event_type = undefined
  dateRange.value = null
  currentPage.value = 1
  loadLogs()
}

function handleDateChange(_val: [string, string] | null) {
  // Date range change triggers search automatically via @change
}

function handleSizeChange() {
  currentPage.value = 1
  loadLogs()
}

function handlePageChange() {
  loadLogs()
}

function getEventTypeLabel(eventType: string): string {
  const found = EVENT_TYPES.find(t => t.value === eventType)
  return found ? found.label : eventType
}

function getEventTypeClass(eventType: string): string {
  const classes: Record<string, string> = {
    login: 'success',
    logout: 'info',
    permission_change: 'warning',
    access_denied: 'danger',
    project_created: 'success',
    project_deleted: 'danger',
    project_archived: 'warning',
    member_added: 'success',
    member_removed: 'warning',
    member_role_changed: 'warning'
  }
  return classes[eventType] || 'default'
}

function formatDateTime(dateStr: string): string {
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}
</script>

<style scoped>
.admin-page {
  padding: var(--space-6);
  max-width: 1400px;
  margin: 0 auto;
}

/* Page Header */
.page-header {
  margin-bottom: var(--space-6);
}

.header-content {
  margin-bottom: var(--space-5);
}

.page-title-section h1 {
  font-family: var(--font-display);
  font-size: var(--font-size-2xl);
  font-weight: var(--font-weight-bold);
  color: var(--color-text-primary);
  margin: 0 0 var(--space-1);
}

.page-title-section p {
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
  margin: 0;
}

/* Filter Form */
.filter-form {
  background: var(--color-bg-surface);
  padding: var(--space-5);
  border-radius: var(--radius-xl);
  border: 1px solid var(--color-border);
}

.filter-row {
  display: flex;
  gap: var(--space-4);
  align-items: flex-end;
  flex-wrap: wrap;
}

.filter-item {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.filter-item label {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  color: var(--color-text-secondary);
}

.filter-input,
.filter-select {
  padding: var(--space-2) var(--space-3);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  font-size: var(--font-size-sm);
  min-width: 120px;
  transition: all var(--transition-base);
}

.filter-input:focus,
.filter-select:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px rgba(91, 95, 199, 0.1);
}

.date-range {
  flex: 1;
  min-width: 280px;
}

.date-picker {
  width: 100%;
}

.filter-actions {
  display: flex;
  gap: var(--space-2);
}

/* Buttons */
.btn {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-4);
  border: none;
  border-radius: var(--radius-md);
  cursor: pointer;
  font-family: var(--font-body);
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  transition: all var(--transition-base);
}

.btn-sm {
  padding: var(--space-1) var(--space-3);
}

.btn-primary {
  background: var(--color-primary);
  color: white;
}

.btn-primary:hover {
  background: var(--color-primary-hover);
}

.btn-ghost {
  background: transparent;
  color: var(--color-text-secondary);
  border: 1px solid var(--color-border);
}

.btn-ghost:hover {
  background: var(--color-bg-muted);
  color: var(--color-text-primary);
}

/* Loading/Error */
.loading-state,
.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--space-16);
  gap: var(--space-4);
  color: var(--color-text-muted);
}

.error-state {
  color: var(--color-danger);
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 3px solid var(--color-border);
  border-top-color: var(--color-primary);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.error-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 64px;
  height: 64px;
  background: rgba(239, 71, 111, 0.1);
  border-radius: 50%;
}

/* Logs Table */
.logs-table-container {
  background: var(--color-bg-surface);
  border-radius: var(--radius-xl);
  border: 1px solid var(--color-border);
  overflow: hidden;
}

.logs-table {
  width: 100%;
  border-collapse: collapse;
}

.logs-table th,
.logs-table td {
  padding: var(--space-3) var(--space-4);
  text-align: left;
  border-bottom: 1px solid var(--color-border);
}

.logs-table th {
  background: var(--color-bg-muted);
  font-weight: var(--font-weight-semibold);
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
}

.log-row:hover {
  background: var(--color-bg-muted);
}

.log-row:last-child td {
  border-bottom: none;
}

.id-cell {
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
  width: 60px;
}

.user-email {
  font-weight: var(--font-weight-medium);
  color: var(--color-text-primary);
}

.event-tag {
  display: inline-block;
  padding: 2px var(--space-2);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-medium);
}

.event-tag.success {
  background: rgba(78, 205, 196, 0.15);
  color: #2e8b7d;
}

.event-tag.info {
  background: var(--color-bg-muted);
  color: var(--color-text-secondary);
}

.event-tag.warning {
  background: rgba(255, 209, 102, 0.15);
  color: #b8860b;
}

.event-tag.danger {
  background: rgba(239, 71, 111, 0.1);
  color: var(--color-danger);
}

.event-tag.default {
  background: rgba(91, 95, 199, 0.1);
  color: var(--color-primary);
}

.details-cell {
  max-width: 200px;
}

.details-text {
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
}

.ip-cell {
  font-size: var(--font-size-sm);
  color: var(--color-text-muted);
  font-family: monospace;
}

.time-cell {
  font-size: var(--font-size-sm);
  color: var(--color-text-muted);
  white-space: nowrap;
}

/* Pagination */
.pagination-wrapper {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: var(--space-4);
  padding: var(--space-3) var(--space-4);
  background: var(--color-bg-surface);
  border-radius: var(--radius-lg);
  border: 1px solid var(--color-border);
}

.pagination-info {
  font-size: var(--font-size-sm);
  color: var(--color-text-muted);
}

/* Responsive */
@media (max-width: 1024px) {
  .filter-row {
    flex-direction: column;
    align-items: stretch;
  }

  .filter-item {
    width: 100%;
  }

  .filter-input,
  .filter-select,
  .date-range {
    width: 100%;
    min-width: unset;
  }

  .filter-actions {
    width: 100%;
    justify-content: flex-end;
  }
}

@media (max-width: 768px) {
  .admin-page {
    padding: var(--space-4);
  }

  .logs-table-container {
    overflow-x: auto;
  }
}
</style>
