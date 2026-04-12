<template>
  <div class="token-manager">
    <!-- Token Created Dialog -->
    <div v-if="showTokenDialog" class="modal-overlay" @click.self="closeTokenDialog">
      <div class="modal-content token-result-modal">
        <div class="modal-header">
          <h3>令牌创建成功</h3>
          <button class="close-btn" @click="closeTokenDialog">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"/>
              <line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>
        <div class="modal-body">
          <div class="warning-box">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/>
              <line x1="12" y1="9" x2="12" y2="13"/>
              <line x1="12" y1="17" x2="12.01" y2="17"/>
            </svg>
            <span>请立即复制并保存您的令牌，它不会再显示第二次！</span>
          </div>
          <div class="token-display">
            <code>{{ newToken }}</code>
            <button class="copy-btn" @click="copyToken">
              <svg v-if="!copied" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="9" y="9" width="13" height="13" rx="2" ry="2"/>
                <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
              </svg>
              <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="20 6 9 17 4 12"/>
              </svg>
              {{ copied ? '已复制' : '复制' }}
            </button>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-primary" @click="closeTokenDialog">关闭</button>
        </div>
      </div>
    </div>

    <!-- Revoke Confirm Dialog -->
    <div v-if="showRevokeDialog" class="modal-overlay" @click.self="closeRevokeDialog">
      <div class="modal-content">
        <div class="modal-header">
          <h3>确认撤销令牌</h3>
          <button class="close-btn" @click="closeRevokeDialog">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"/>
              <line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>
        <div class="modal-body">
          <p>确定要撤销令牌 <strong>{{ revokingToken?.name }}</strong> 吗？此操作不可撤销。</p>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="closeRevokeDialog">取消</button>
          <button class="btn btn-danger" @click="confirmRevoke" :disabled="revoking">
            {{ revoking ? '撤销中...' : '确认撤销' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Create Token Form -->
    <div class="section-card">
      <div class="section-header">
        <h3>创建新令牌</h3>
        <p>生成一个新的个人访问令牌</p>
      </div>
      <form @submit.prevent="handleCreateToken" class="token-form">
        <div class="form-group">
          <label for="token_name">令牌名称</label>
          <input
            id="token_name"
            v-model="createForm.name"
            type="text"
            placeholder="例如：开发环境、MCP客户端"
            required
            maxlength="100"
            class="form-input"
          />
        </div>
        <div class="form-group">
          <label for="token_scope">权限范围</label>
          <select
            id="token_scope"
            v-model="createForm.scope"
            required
            class="form-input"
          >
            <option value="">选择权限范围</option>
            <option value="task:read">task:read - 读取任务</option>
            <option value="task:read,task:write">task:read,task:write - 读写任务</option>
            <option value="task:read,project:read,comment:read">task:read,project:read,comment:read - 只读所有</option>
            <option value="task:read,task:write,project:read,project:write,comment:read,comment:write">全部权限</option>
          </select>
        </div>
        <div class="form-group">
          <label for="expires_days">过期天数</label>
          <select
            id="expires_days"
            v-model.number="createForm.expiresInDays"
            class="form-input"
          >
            <option :value="30">30 天</option>
            <option :value="60">60 天</option>
            <option :value="90">90 天</option>
            <option :value="180">180 天</option>
            <option :value="365">365 天</option>
          </select>
        </div>
        <div class="form-actions">
          <button type="submit" class="btn btn-primary" :disabled="creating">
            <svg v-if="!creating" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="12" y1="5" x2="12" y2="19"/>
              <line x1="5" y1="12" x2="19" y2="12"/>
            </svg>
            {{ creating ? '创建中...' : '创建令牌' }}
          </button>
        </div>
      </form>
      <Transition name="fade">
        <div v-if="createError" class="error-message">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/>
            <line x1="12" y1="8" x2="12" y2="12"/>
            <line x1="12" y1="16" x2="12.01" y2="16"/>
          </svg>
          {{ createError }}
        </div>
      </Transition>
    </div>

    <!-- Token List -->
    <div class="section-card">
      <div class="section-header">
        <h3>我的令牌</h3>
        <p>已创建的个人访问令牌</p>
      </div>

      <div v-if="loading" class="loading-state">
        <div class="loading-spinner"></div>
        <span>加载中...</span>
      </div>

      <div v-else-if="tokens.length === 0" class="empty-state">
        <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
          <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
          <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
        </svg>
        <p>暂无令牌</p>
        <span>创建您的第一个个人访问令牌</span>
      </div>

      <div v-else class="token-list">
        <div
          v-for="token in tokens"
          :key="token.id"
          :class="['token-item', { revoked: token.revoked }]"
        >
          <div class="token-info">
            <div class="token-name-row">
              <span class="token-name">{{ token.name }}</span>
              <span v-if="token.revoked" class="revoked-badge">已撤销</span>
              <span v-else-if="isExpired(token.expires_at)" class="expired-badge">已过期</span>
            </div>
            <div class="token-meta">
              <span class="token-prefix">{{ token.token_prefix }}...</span>
              <span class="separator">·</span>
              <span class="token-scope">{{ token.scope }}</span>
            </div>
            <div class="token-dates">
              <span>创建于 {{ formatDate(token.created_at) }}</span>
              <span v-if="token.last_used_at">，最后使用于 {{ formatDate(token.last_used_at) }}</span>
            </div>
          </div>
          <div class="token-actions" v-if="!token.revoked">
            <button class="btn-icon revoke-btn" @click="openRevokeDialog(token)" title="撤销令牌">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="3 6 5 6 21 6"/>
                <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
              </svg>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import { useTokenStore } from '../../stores/token'
import type { TokenInfo } from '../../api/token'

const tokenStore = useTokenStore()
const { tokens, loading, error } = storeToRefs(tokenStore)

const creating = ref(false)
const createError = ref('')
const showTokenDialog = ref(false)
const newToken = ref('')
const copied = ref(false)

const showRevokeDialog = ref(false)
const revokingToken = ref<TokenInfo | null>(null)
const revoking = ref(false)

const createForm = reactive({
  name: '',
  scope: '',
  expiresInDays: 90
})

onMounted(async () => {
  await fetchTokens()
})

async function fetchTokens() {
  loading.value = true
  await tokenStore.fetchTokens()
  loading.value = false
}

async function handleCreateToken() {
  createError.value = ''
  creating.value = true

  try {
    const result = await tokenStore.createToken({
      name: createForm.name,
      scope: createForm.scope,
      expires_in_days: createForm.expiresInDays
    })

    if (result.success && result.token) {
      newToken.value = result.token
      showTokenDialog.value = true
      createForm.name = ''
      createForm.scope = ''
      createForm.expiresInDays = 90
    } else {
      createError.value = error.value || '创建失败'
    }
  } finally {
    creating.value = false
  }
}

function closeTokenDialog() {
  showTokenDialog.value = false
  newToken.value = ''
  copied.value = false
  tokenStore.clearLastCreatedToken()
}

async function copyToken() {
  try {
    await navigator.clipboard.writeText(newToken.value)
    copied.value = true
    setTimeout(() => { copied.value = false }, 2000)
  } catch (e) {
    console.error('复制失败', e)
  }
}

function openRevokeDialog(token: TokenInfo) {
  revokingToken.value = token
  showRevokeDialog.value = true
}

function closeRevokeDialog() {
  showRevokeDialog.value = false
  revokingToken.value = null
}

async function confirmRevoke() {
  if (!revokingToken.value) return

  revoking.value = true
  try {
    const success = await tokenStore.revokeToken(revokingToken.value.id)
    if (success) {
      closeRevokeDialog()
    }
  } finally {
    revoking.value = false
  }
}

function isExpired(expiresAt: string): boolean {
  return new Date(expiresAt) < new Date()
}

function formatDate(dateStr: string): string {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}
</script>

<style scoped>
.token-manager {
  display: flex;
  flex-direction: column;
  gap: var(--space-6);
}

/* Section Card */
.section-card {
  background: var(--color-bg-surface);
  border-radius: var(--radius-xl);
  padding: var(--space-6);
  box-shadow: var(--shadow-card);
  border: 1px solid var(--color-border);
}

.section-header {
  margin-bottom: var(--space-6);
  padding-bottom: var(--space-4);
  border-bottom: 1px solid var(--color-border);
}

.section-header h3 {
  font-family: var(--font-display);
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-primary);
  margin: 0 0 var(--space-1);
}

.section-header p {
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
  margin: 0;
}

/* Form */
.token-form {
  display: flex;
  flex-direction: column;
  gap: var(--space-5);
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.form-group label {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  color: var(--color-text-secondary);
}

.form-input {
  padding: var(--space-3) var(--space-4);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  font-family: var(--font-body);
  font-size: var(--font-size-base);
  color: var(--color-text-primary);
  background: var(--color-bg-surface);
  transition: all var(--transition-base);
}

.form-input:hover {
  border-color: var(--color-border-hover);
}

.form-input:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px rgba(91, 95, 199, 0.1);
}

.form-actions {
  display: flex;
  justify-content: flex-end;
}

/* Buttons */
.btn {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-3) var(--space-5);
  border: none;
  border-radius: var(--radius-lg);
  cursor: pointer;
  font-family: var(--font-body);
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  transition: all var(--transition-base);
}

.btn-primary {
  background: var(--color-primary);
  color: white;
  box-shadow: 0 2px 8px rgba(91, 95, 199, 0.3);
}

.btn-primary:hover:not(:disabled) {
  background: var(--color-primary-hover);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(91, 95, 199, 0.4);
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-secondary {
  background: var(--color-bg-base);
  color: var(--color-text-primary);
  border: 1px solid var(--color-border);
}

.btn-secondary:hover {
  background: var(--color-bg-hover);
}

.btn-danger {
  background: var(--color-danger);
  color: white;
}

.btn-danger:hover:not(:disabled) {
  opacity: 0.9;
}

.btn-danger:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border: none;
  border-radius: var(--radius-lg);
  cursor: pointer;
  transition: all var(--transition-base);
  background: transparent;
}

.revoke-btn {
  color: var(--color-text-muted);
}

.revoke-btn:hover {
  background: rgba(239, 71, 111, 0.1);
  color: var(--color-danger);
}

/* Token List */
.token-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

.token-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space-4);
  background: var(--color-bg-base);
  border-radius: var(--radius-lg);
  border: 1px solid var(--color-border);
  transition: all var(--transition-base);
}

.token-item:hover {
  border-color: var(--color-border-hover);
}

.token-item.revoked {
  opacity: 0.6;
}

.token-info {
  flex: 1;
}

.token-name-row {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  margin-bottom: var(--space-1);
}

.token-name {
  font-weight: var(--font-weight-medium);
  color: var(--color-text-primary);
}

.revoked-badge,
.expired-badge {
  font-size: var(--font-size-xs);
  padding: 2px 8px;
  border-radius: var(--radius-full);
  font-weight: var(--font-weight-medium);
}

.revoked-badge {
  background: rgba(239, 71, 111, 0.15);
  color: var(--color-danger);
}

.expired-badge {
  background: rgba(255, 209, 102, 0.15);
  color: #b8860b;
}

.token-meta {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
  margin-bottom: var(--space-1);
}

.token-prefix {
  font-family: monospace;
  background: var(--color-bg-surface);
  padding: 2px 6px;
  border-radius: var(--radius-sm);
}

.separator {
  color: var(--color-text-muted);
}

.token-scope {
  color: var(--color-primary);
}

.token-dates {
  font-size: var(--font-size-xs);
  color: var(--color-text-muted);
}

/* Loading/Empty States */
.loading-state,
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--space-12);
  color: var(--color-text-muted);
  text-align: center;
}

.loading-state {
  gap: var(--space-3);
}

.empty-state {
  gap: var(--space-3);
}

.empty-state svg {
  opacity: 0.5;
}

.empty-state p {
  margin: 0;
  font-weight: var(--font-weight-medium);
  color: var(--color-text-secondary);
}

.empty-state span {
  font-size: var(--font-size-sm);
}

.loading-spinner {
  width: 32px;
  height: 32px;
  border: 3px solid var(--color-border);
  border-top-color: var(--color-primary);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* Modal */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: var(--space-4);
}

.modal-content {
  background: var(--color-bg-surface);
  border-radius: var(--radius-xl);
  width: 100%;
  max-width: 480px;
  box-shadow: var(--shadow-lg);
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space-5) var(--space-6);
  border-bottom: 1px solid var(--color-border);
}

.modal-header h3 {
  font-family: var(--font-display);
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-primary);
  margin: 0;
}

.close-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: none;
  border-radius: var(--radius-md);
  background: transparent;
  color: var(--color-text-muted);
  cursor: pointer;
  transition: all var(--transition-base);
}

.close-btn:hover {
  background: var(--color-bg-hover);
  color: var(--color-text-primary);
}

.modal-body {
  padding: var(--space-6);
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: var(--space-3);
  padding: var(--space-4) var(--space-6);
  border-top: 1px solid var(--color-border);
}

/* Token Result Modal */
.token-result-modal .modal-body {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
}

.warning-box {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-4);
  background: rgba(255, 209, 102, 0.15);
  border-radius: var(--radius-lg);
  color: #b8860b;
  font-size: var(--font-size-sm);
}

.warning-box svg {
  flex-shrink: 0;
}

.token-display {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-4);
  background: var(--color-bg-base);
  border-radius: var(--radius-lg);
  border: 1px solid var(--color-border);
}

.token-display code {
  flex: 1;
  font-family: monospace;
  font-size: var(--font-size-sm);
  word-break: break-all;
  color: var(--color-text-primary);
}

.copy-btn {
  display: inline-flex;
  align-items: center;
  gap: var(--space-1);
  padding: var(--space-2) var(--space-3);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  background: var(--color-bg-surface);
  color: var(--color-text-secondary);
  font-size: var(--font-size-xs);
  cursor: pointer;
  transition: all var(--transition-base);
  white-space: nowrap;
}

.copy-btn:hover {
  background: var(--color-bg-hover);
  color: var(--color-text-primary);
}

/* Messages */
.error-message {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  margin-top: var(--space-4);
  padding: var(--space-3) var(--space-4);
  border-radius: var(--radius-lg);
  background: rgba(239, 71, 111, 0.1);
  color: var(--color-danger);
  font-size: var(--font-size-sm);
}

/* Fade Animation */
.fade-enter-active,
.fade-leave-active {
  transition: opacity var(--transition-base);
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>