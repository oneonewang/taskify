<template>
  <div class="team-members">
    <div class="members-header">
      <span class="header-title">团队成员</span>
      <span class="member-count">{{ users.length }}</span>
    </div>
    <div class="members-list">
      <div
        v-for="user in users"
        :key="user.id"
        class="member-item"
        :class="{ 'is-current': user.id === currentUserId }"
      >
        <div class="member-avatar" :style="{ backgroundColor: user.avatar }">
          {{ user.name.charAt(0) }}
        </div>
        <div class="member-info">
          <div class="member-name">{{ user.name }}</div>
          <div class="member-role">{{ roleLabel(user.role) }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getUsers, type User } from '../api/user'

defineProps<{
  currentUserId: number | null
}>()

const users = ref<User[]>([])
const loading = ref(true)

onMounted(async () => {
  try {
    const res = await getUsers()
    if (res.code === 0) {
      users.value = res.data
    }
  } catch (e) {
    console.error('Failed to load users:', e)
  } finally {
    loading.value = false
  }
})

function roleLabel(role: string): string {
  return role === 'product_manager' ? '产品经理' : '工程师'
}
</script>

<style scoped>
.team-members {
  background: white;
  border-radius: 12px;
  padding: 16px;
  min-width: 240px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.members-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid #eee;
}

.header-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
}

.member-count {
  background: #667eea;
  color: white;
  padding: 2px 8px;
  border-radius: 10px;
  font-size: 12px;
}

.members-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.member-item {
  display: flex;
  align-items: center;
  padding: 8px;
  border-radius: 8px;
  transition: background 0.2s;
}

.member-item:hover {
  background: #f5f7fa;
}

.member-item.is-current {
  background: #f0f1ff;
  border: 1px solid #667eea;
}

.member-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 14px;
  font-weight: 600;
  margin-right: 12px;
}

.member-info {
  flex: 1;
}

.member-name {
  font-size: 14px;
  font-weight: 500;
  color: #333;
}

.member-role {
  font-size: 12px;
  color: #999;
  margin-top: 2px;
}
</style>