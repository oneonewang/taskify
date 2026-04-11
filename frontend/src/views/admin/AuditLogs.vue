<template>
  <div class="audit-logs-page">
    <div class="page-header">
      <h1>审计日志</h1>
    </div>

    <!-- 筛选表单 -->
    <div class="filter-form">
      <el-form :inline="true" :model="filterForm">
        <el-form-item label="用户">
          <el-input
            v-model="filterForm.user_id"
            type="number"
            placeholder="用户ID"
            clearable
            @keyup.enter="handleSearch"
          />
        </el-form-item>
        <el-form-item label="事件类型">
          <el-select
            v-model="filterForm.event_type"
            placeholder="选择事件类型"
            clearable
            @change="handleSearch"
          >
            <el-option
              v-for="type in EVENT_TYPES"
              :key="type.value"
              :label="type.label"
              :value="type.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="时间范围">
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
            @change="handleDateChange"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 日志列表 -->
    <div v-if="loading" class="loading">加载中...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <template v-else>
      <el-table :data="logs" stripe style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="user_email" label="用户邮箱" min-width="180">
          <template #default="{ row }">
            {{ row.user_email || `用户ID: ${row.user_id}` }}
          </template>
        </el-table-column>
        <el-table-column prop="event_type" label="事件类型" width="140">
          <template #default="{ row }">
            <el-tag :type="getEventTypeTag(row.event_type)">
              {{ getEventTypeLabel(row.event_type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="details" label="详情" min-width="200" show-overflow-tooltip />
        <el-table-column prop="ip_address" label="IP地址" width="140" />
        <el-table-column prop="created_at" label="时间" width="180">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
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

function handleDateChange(val: [string, string] | null) {
  if (!val) {
    filterForm.user_id = undefined
  }
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

function getEventTypeTag(eventType: string): string {
  const tags: Record<string, string> = {
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
  return tags[eventType] || 'info'
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
.audit-logs-page {
  padding: 24px;
}

.page-header {
  margin-bottom: 24px;
}

.page-header h1 {
  margin: 0;
  font-size: 24px;
  color: #333;
}

.filter-form {
  background: #fff;
  padding: 16px;
  border-radius: 8px;
  margin-bottom: 16px;
}

.loading,
.error {
  text-align: center;
  padding: 40px;
  color: #666;
}

.error {
  color: #dc3545;
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>
