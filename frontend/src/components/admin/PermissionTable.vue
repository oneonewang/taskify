<template>
  <div class="permission-table">
    <div v-if="loading" class="loading">加载权限中...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <template v-else>
      <div v-for="group in permissionGroups" :key="group.resource" class="permission-group">
        <div class="group-header">
          <h4>{{ getResourceDisplayName(group.resource) }}</h4>
        </div>
        <div class="permission-list">
          <div
            v-for="perm in group.permissions"
            :key="perm.id"
            class="permission-item"
            :class="{ 'is-selected': selectedPermissions.has(perm.id) }"
            @click="togglePermission(perm.id)"
          >
            <el-checkbox
              :model-value="selectedPermissions.has(perm.id)"
              @change="() => togglePermission(perm.id)"
            />
            <div class="permission-info">
              <span class="permission-action">{{ perm.action }}</span>
              <span class="permission-desc">{{ perm.description || '' }}</span>
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { getPermissions, type Permission } from '../../api/role'

const props = defineProps<{
  modelValue: number[]
}>()

const emit = defineEmits<{
  'update:modelValue': [value: number[]]
}>()

const permissions = ref<Permission[]>([])
const loading = ref(false)
const error = ref<string | null>(null)

const selectedPermissions = computed({
  get: () => new Set(props.modelValue),
  set: (val) => emit('update:modelValue', Array.from(val))
})

onMounted(async () => {
  await loadPermissions()
})

async function loadPermissions() {
  loading.value = true
  error.value = null
  try {
    const res = await getPermissions()
    if (res.code === 0) {
      permissions.value = res.data.permissions
    } else {
      error.value = res.message
    }
  } catch (e: any) {
    error.value = e.message || '加载权限失败'
  } finally {
    loading.value = false
  }
}

const permissionGroups = computed(() => {
  const groups: { resource: string; permissions: Permission[] }[] = []
  const resourceMap = new Map<string, Permission[]>()

  for (const perm of permissions.value) {
    const list = resourceMap.get(perm.resource) || []
    list.push(perm)
    resourceMap.set(perm.resource, list)
  }

  for (const [resource, perms] of resourceMap) {
    groups.push({ resource, permissions: perms })
  }

  return groups
})

function getResourceDisplayName(resource: string): string {
  const names: Record<string, string> = {
    users: '用户管理',
    projects: '项目管理',
    tasks: '任务管理',
    comments: '评论管理',
    roles: '角色管理'
  }
  return names[resource] || resource
}

function togglePermission(permId: number) {
  const newSet = new Set(selectedPermissions.value)
  if (newSet.has(permId)) {
    newSet.delete(permId)
  } else {
    newSet.add(permId)
  }
  selectedPermissions.value = newSet
}
</script>

<style scoped>
.permission-table {
  max-height: 400px;
  overflow-y: auto;
}

.loading,
.error {
  text-align: center;
  padding: 20px;
  color: #666;
}

.error {
  color: #dc3545;
}

.permission-group {
  margin-bottom: 16px;
}

.group-header {
  margin-bottom: 8px;
}

.group-header h4 {
  margin: 0;
  font-size: 14px;
  color: #333;
  font-weight: 600;
  text-transform: capitalize;
}

.permission-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding-left: 12px;
}

.permission-item {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  padding: 6px 8px;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.permission-item:hover {
  background-color: #f5f7fa;
}

.permission-item.is-selected {
  background-color: #e3f2fd;
}

.permission-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.permission-action {
  font-size: 13px;
  font-weight: 500;
  color: #333;
  text-transform: capitalize;
}

.permission-desc {
  font-size: 12px;
  color: #999;
}
</style>
