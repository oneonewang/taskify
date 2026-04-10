<template>
  <el-dialog
    v-model="dialogVisible"
    :title="isEdit ? '编辑任务' : '创建任务'"
    width="500px"
    @close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="formData"
      :rules="rules"
      label-width="80px"
    >
      <el-form-item label="标题" prop="title">
        <el-input
          v-model="formData.title"
          placeholder="请输入任务标题"
          maxlength="200"
          show-word-limit
        />
      </el-form-item>

      <el-form-item label="描述" prop="description">
        <el-input
          v-model="formData.description"
          type="textarea"
          placeholder="请输入任务描述（可选）"
          :rows="3"
          maxlength="2000"
          show-word-limit
        />
      </el-form-item>

      <el-form-item label="负责人" prop="assignee_id">
        <el-select
          v-model="formData.assignee_id"
          placeholder="请选择负责人"
          style="width: 100%"
        >
          <el-option
            v-for="user in assigneeOptions"
            :key="user.id"
            :label="user.name"
            :value="user.id"
          >
            <div class="user-option">
              <span
                class="user-avatar"
                :style="{ backgroundColor: user.avatar }"
              >
                {{ user.name.charAt(0) }}
              </span>
              <span>{{ user.name }}</span>
            </div>
          </el-option>
        </el-select>
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button type="primary" :loading="loading" @click="handleSubmit">
        {{ isEdit ? '保存' : '创建' }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage } from 'element-plus'
import { createTask, updateTask, type CreateTaskRequest, type UpdateTaskRequest } from '../../api/task'
import { getProjectMembers } from '../../api/membership'
import type { Task } from '../../api/project'
import type { ProjectMember } from '../../api/membership'

const props = defineProps<{
  visible: boolean
  task?: Task | null
  projectId: number
}>()

const emit = defineEmits<{
  'update:visible': [value: boolean]
  success: [task: Task]
}>()

const dialogVisible = computed({
  get: () => props.visible,
  set: (val) => emit('update:visible', val)
})

const isEdit = computed(() => !!props.task)

const formRef = ref<FormInstance>()
const loading = ref(false)
const users = ref<ProjectMember[]>([])

// 转换为 Assignee 格式供下拉选择
const assigneeOptions = computed(() => users.value.map(m => ({
  id: m.user_id,
  name: m.display_name,
  avatar: m.avatar_url || '#667eea'
})))

// 表单数据
const formData = ref({
  title: '',
  description: '',
  assignee_id: null as number | null
})

// 表单验证规则
const rules: FormRules = {
  title: [
    { required: true, message: '请输入任务标题', trigger: 'blur' },
    { max: 200, message: '标题最多200个字符', trigger: 'blur' }
  ],
  assignee_id: [
    { required: true, message: '请选择负责人', trigger: 'change' }
  ]
}

// 加载项目成员列表
async function loadUsers() {
  if (!props.projectId) return
  try {
    const res = await getProjectMembers(props.projectId)
    if (res.code === 0 && res.data.members) {
      users.value = res.data.members
    }
  } catch (e) {
    console.error('Failed to load members:', e)
  }
}

// 监听visible变化，初始化表单
watch(() => props.visible, async (val) => {
  if (val) {
    await loadUsers()
    if (props.task) {
      // 编辑模式
      formData.value = {
        title: props.task.title,
        description: props.task.description || '',
        assignee_id: props.task.assignee.id
      }
    } else {
      // 创建模式
      formData.value = {
        title: '',
        description: '',
        assignee_id: null
      }
    }
  }
})

// 提交表单
async function handleSubmit() {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
  } catch {
    return
  }

  loading.value = true

  try {
    if (isEdit.value && props.task) {
      // 更新任务
      const updateData: UpdateTaskRequest = {}
      if (formData.value.title !== props.task.title) {
        updateData.title = formData.value.title
      }
      if (formData.value.description !== (props.task.description || '')) {
        updateData.description = formData.value.description
      }
      if (formData.value.assignee_id !== props.task.assignee.id) {
        updateData.assignee_id = formData.value.assignee_id!
      }

      const res = await updateTask(props.task.id, updateData)
      if (res.code === 0) {
        ElMessage.success('任务更新成功')
        emit('success', res.data)
        handleClose()
      } else {
        ElMessage.error(res.message || '更新任务失败')
      }
    } else {
      // 创建任务
      const createData: CreateTaskRequest = {
        title: formData.value.title,
        description: formData.value.description || undefined,
        assignee_id: formData.value.assignee_id!
      }

      const res = await createTask(props.projectId, createData)
      if (res.code === 0) {
        ElMessage.success('任务创建成功')
        emit('success', res.data)
        handleClose()
      } else {
        ElMessage.error(res.message || '创建任务失败')
      }
    }
  } catch (e: any) {
    ElMessage.error(e.message || '操作失败')
  } finally {
    loading.value = false
  }
}

// 关闭对话框
function handleClose() {
  formRef.value?.resetFields()
  dialogVisible.value = false
}
</script>

<style scoped>
.user-option {
  display: flex;
  align-items: center;
  gap: 8px;
}

.user-avatar {
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
</style>
