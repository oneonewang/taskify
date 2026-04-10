<template>
  <el-dialog
    v-model="dialogVisible"
    title="任务详情"
    width="600px"
    @close="handleClose"
  >
    <template v-if="task">
    <!-- 任务基本信息 -->
    <div class="task-info">
      <h2 class="task-title">{{ task.title }}</h2>
      <div class="task-meta">
        <el-tag :type="statusType" size="small">{{ statusText }}</el-tag>
        <span class="assignee">
          <span
            class="assignee-avatar"
            :style="{ backgroundColor: task.assignee.avatar }"
          >
            {{ task.assignee.name.charAt(0) }}
          </span>
          {{ task.assignee.name }}
        </span>
      </div>
      <div v-if="task.description" class="task-description">
        {{ task.description }}
      </div>
    </div>

    <el-divider />

    <!-- 评论区域 -->
    <div class="comments-section">
      <h3 class="comments-title">评论 ({{ comments.length }})</h3>

      <!-- 评论列表 -->
      <div class="comments-list">
        <div
          v-for="comment in comments"
          :key="comment.id"
          class="comment-item"
        >
          <div class="comment-header">
            <div class="comment-author">
              <span
                class="comment-avatar"
                :style="{ backgroundColor: comment.user.avatar }"
              >
                {{ comment.user.name.charAt(0) }}
              </span>
              <span class="comment-name">{{ comment.user.name }}</span>
              <span class="comment-time">{{ formatTime(comment.created_at) }}</span>
            </div>
            <div class="comment-actions">
              <!-- 编辑按钮：仅当前用户的评论显示 -->
              <template v-if="isCurrentUser(comment.user.id)">
                <el-button
                  type="primary"
                  text
                  size="small"
                  @click="startEdit(comment)"
                >
                  编辑
                </el-button>
                <el-button
                  type="danger"
                  text
                  size="small"
                  @click="confirmDelete(comment)"
                >
                  删除
                </el-button>
              </template>
            </div>
          </div>

          <!-- 评论内容或编辑表单 -->
          <div v-if="editingCommentId !== comment.id" class="comment-content">
            {{ comment.content }}
          </div>
          <div v-else class="comment-edit-form">
            <el-input
              v-model="editContent"
              type="textarea"
              :rows="2"
              placeholder="请输入评论内容"
              maxlength="2000"
              show-word-limit
            />
            <div class="edit-actions">
              <el-button size="small" @click="cancelEdit">取消</el-button>
              <el-button
                type="primary"
                size="small"
                :loading="editLoading"
                @click="submitEdit(comment.id)"
              >
                保存
              </el-button>
            </div>
          </div>
        </div>

        <div v-if="comments.length === 0 && !loading" class="no-comments">
          暂无评论，添加第一条评论吧
        </div>
      </div>

      <!-- 添加评论表单 -->
      <div class="add-comment">
        <el-input
          v-model="newComment"
          type="textarea"
          :rows="3"
          placeholder="添加评论..."
          maxlength="2000"
          show-word-limit
        />
        <el-button
          type="primary"
          :loading="submitLoading"
          :disabled="!newComment.trim()"
          @click="submitComment"
        >
          提交评论
        </el-button>
      </div>
    </div>
    </template>
  </el-dialog>

  <!-- 删除确认对话框 -->
  <el-dialog
    v-model="deleteDialogVisible"
    title="确认删除"
    width="400px"
    append-to-body
  >
    <p>确定要删除这条评论吗？此操作无法撤销。</p>
    <template #footer>
      <el-button @click="deleteDialogVisible = false">取消</el-button>
      <el-button type="danger" :loading="deleteLoading" @click="submitDelete">
        删除
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import type { Task } from '../../api/project'
import { getComments, createComment, updateComment, deleteComment, type Comment } from '../../api/comment'
import { useAuthStore } from '../../stores/auth'
import { useSSE } from '../../composables/useSSE'

const props = defineProps<{
  visible: boolean
  task: Task | null
}>()

const emit = defineEmits<{
  'update:visible': [value: boolean]
}>()

const authStore = useAuthStore()

// SSE 监听评论添加事件
useSSE({
  onCommentAdded: (data) => {
    // 仅当弹窗打开且评论属于当前任务时更新
    if (props.visible && props.task && data.task_id === props.task.id) {
      // 检查是否已存在（防止重复）
      const exists = comments.value.some(c => c.id === data.id)
      if (!exists) {
        comments.value.push(data)
        ElMessage.info('收到新评论')
      }
    }
  }
})

const dialogVisible = computed({
  get: () => props.visible,
  set: (val) => emit('update:visible', val)
})

// 状态映射
const statusMap: Record<string, { text: string; type: string }> = {
  todo: { text: '待办', type: 'info' },
  in_progress: { text: '进行中', type: 'warning' },
  review: { text: '审核中', type: 'primary' },
  done: { text: '已完成', type: 'success' }
}

const statusText = computed(() => statusMap[props.task?.status || '']?.text || props.task?.status)
const statusType = computed(() => statusMap[props.task?.status || '']?.type || 'info')

// 评论相关
const comments = ref<Comment[]>([])
const loading = ref(false)
const newComment = ref('')
const submitLoading = ref(false)

// 编辑相关
const editingCommentId = ref<number | null>(null)
const editContent = ref('')
const editLoading = ref(false)

// 删除相关
const deleteDialogVisible = ref(false)
const deleteLoading = ref(false)
const commentToDelete = ref<Comment | null>(null)

// 判断是否是当前用户
function isCurrentUser(userId: number): boolean {
  return authStore.user?.id === userId
}

// 格式化时间
function formatTime(timeStr: string): string {
  const date = new Date(timeStr)
  const now = new Date()
  const diff = now.getTime() - date.getTime()

  if (diff < 60000) return '刚刚'
  if (diff < 3600000) return `${Math.floor(diff / 60000)}分钟前`
  if (diff < 86400000) return `${Math.floor(diff / 3600000)}小时前`

  if (date.toDateString() === now.toDateString()) {
    return date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  }

  return date.toLocaleDateString('zh-CN', { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })
}

// 加载评论
async function loadComments() {
  if (!props.task) return

  loading.value = true
  try {
    const res = await getComments(props.task.id)
    if (res.code === 0) {
      comments.value = res.data
    }
  } catch (e) {
    console.error('Failed to load comments:', e)
  } finally {
    loading.value = false
  }
}

// 提交新评论
async function submitComment() {
  if (!props.task || !newComment.value.trim()) return

  submitLoading.value = true
  try {
    const res = await createComment(props.task.id, newComment.value.trim())
    if (res.code === 0) {
      comments.value.push(res.data)
      newComment.value = ''
      ElMessage.success('评论添加成功')
    } else {
      ElMessage.error(res.message || '添加评论失败')
    }
  } catch (e: any) {
    ElMessage.error(e.message || '添加评论失败')
  } finally {
    submitLoading.value = false
  }
}

// 开始编辑评论
function startEdit(comment: Comment) {
  editingCommentId.value = comment.id
  editContent.value = comment.content
}

// 取消编辑
function cancelEdit() {
  editingCommentId.value = null
  editContent.value = ''
}

// 提交编辑
async function submitEdit(commentId: number) {
  if (!props.task || !editContent.value.trim()) return

  editLoading.value = true
  try {
    const res = await updateComment(props.task.id, commentId, editContent.value.trim())
    if (res.code === 0) {
      // 更新本地评论列表
      const index = comments.value.findIndex(c => c.id === commentId)
      if (index !== -1) {
        comments.value[index] = res.data
      }
      cancelEdit()
      ElMessage.success('评论已更新')
    } else {
      ElMessage.error(res.message || '更新评论失败')
    }
  } catch (e: any) {
    ElMessage.error(e.message || '更新评论失败')
  } finally {
    editLoading.value = false
  }
}

// 确认删除
function confirmDelete(comment: Comment) {
  commentToDelete.value = comment
  deleteDialogVisible.value = true
}

// 提交删除
async function submitDelete() {
  if (!props.task || !commentToDelete.value) return

  deleteLoading.value = true
  try {
    const res = await deleteComment(props.task.id, commentToDelete.value.id)
    if (res.code === 0) {
      comments.value = comments.value.filter(c => c.id !== commentToDelete.value!.id)
      deleteDialogVisible.value = false
      commentToDelete.value = null
      ElMessage.success('评论已删除')
    } else {
      ElMessage.error(res.message || '删除评论失败')
    }
  } catch (e: any) {
    ElMessage.error(e.message || '删除评论失败')
  } finally {
    deleteLoading.value = false
  }
}

// 关闭对话框
function handleClose() {
  comments.value = []
  newComment.value = ''
  editingCommentId.value = null
  editContent.value = ''
  deleteDialogVisible.value = false
  commentToDelete.value = null
}

// 监听 visible 变化
watch(() => props.visible, (val) => {
  if (val) {
    loadComments()
  }
})
</script>

<style scoped>
.task-info {
  margin-bottom: 16px;
}

.task-title {
  font-size: 20px;
  font-weight: 600;
  color: #333;
  margin: 0 0 12px 0;
}

.task-meta {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 12px;
}

.assignee {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #666;
  font-size: 14px;
}

.assignee-avatar {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 12px;
  font-weight: 600;
}

.task-description {
  color: #666;
  font-size: 14px;
  line-height: 1.6;
  background: #f5f7fa;
  padding: 12px;
  border-radius: 6px;
}

.comments-section {
  margin-top: 16px;
}

.comments-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin: 0 0 16px 0;
}

.comments-list {
  max-height: 300px;
  overflow-y: auto;
  margin-bottom: 16px;
}

.comment-item {
  padding: 12px 0;
  border-bottom: 1px solid #eee;
}

.comment-item:last-child {
  border-bottom: none;
}

.comment-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.comment-author {
  display: flex;
  align-items: center;
  gap: 8px;
}

.comment-avatar {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 12px;
  font-weight: 600;
}

.comment-name {
  font-size: 14px;
  font-weight: 500;
  color: #333;
}

.comment-time {
  font-size: 12px;
  color: #999;
}

.comment-actions {
  display: flex;
  gap: 8px;
}

.comment-content {
  font-size: 14px;
  color: #666;
  line-height: 1.6;
  padding-left: 36px;
}

.comment-edit-form {
  padding-left: 36px;
}

.edit-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 8px;
}

.no-comments {
  text-align: center;
  color: #999;
  padding: 20px 0;
  font-size: 14px;
}

.add-comment {
  border-top: 1px solid #eee;
  padding-top: 16px;
}

.add-comment .el-input {
  margin-bottom: 12px;
}

.add-comment .el-button {
  float: right;
}
</style>
