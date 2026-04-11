<template>
  <div class="team-members">
    <div class="members-header">
      <span class="header-title">团队成员</span>
      <span class="member-count">{{ members.length }}</span>
    </div>
    <div v-if="loading" class="loading">加载中...</div>
    <div v-else class="members-list">
      <div
        v-for="member in members"
        :key="member.user_id"
        class="member-item"
        :class="{ 'is-current': member.user_id === currentUserId }"
      >
        <div class="member-avatar" :style="{ backgroundColor: member.avatar_url || '#667eea' }">
          {{ member.display_name.charAt(0) }}
        </div>
        <div class="member-info">
          <div class="member-name">{{ member.display_name }}</div>
          <div class="member-role">{{ member.role_display_name || member.role_name }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getProjectMembers } from '../api/membership'

defineProps<{
  currentUserId: number | null
  projectId: number
}>()

interface Member {
  user_id: number
  display_name: string
  avatar_url: string
  role_name: string
  role_display_name: string
}

const members = ref<Member[]>([])
const loading = ref(true)

onMounted(async () => {
  // props 将在父组件确保 projectId 有效后才渲染此组件
})

function loadMembers(projectId: number) {
  loading.value = true
  members.value = []
  getProjectMembers(projectId).then(res => {
    if (res.code === 0 && res.data.members) {
      members.value = res.data.members
      console.log('[DEBUG] TeamMembers loaded:', members.value.length, 'members')
    }
  }).catch(e => {
    console.error('[ERROR] Failed to load members:', e)
    // 不阻塞流程，允许空列表显示
  }).finally(() => {
    loading.value = false
  })
}

// 暴露方法给父组件调用
defineExpose({ loadMembers })
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