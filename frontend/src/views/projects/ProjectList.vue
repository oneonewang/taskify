<template>
  <div class="project-list-page">
    <AppHeader :title="activeTab === 'my' ? '我参与的项目' : '全部项目'" />

    <main class="page-content">
      <!-- Tab Bar -->
      <div class="tab-bar animate-in">
        <button
          :class="['tab-btn', { active: activeTab === 'my' }]"
          @click="switchTab('my')"
        >
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
            <circle cx="12" cy="7" r="4"/>
          </svg>
          我参与的项目
        </button>
        <button
          :class="['tab-btn', { active: activeTab === 'all' }]"
          @click="switchTab('all')"
        >
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="3" y="3" width="7" height="7"/>
            <rect x="14" y="3" width="7" height="7"/>
            <rect x="14" y="14" width="7" height="7"/>
            <rect x="3" y="14" width="7" height="7"/>
          </svg>
          全部项目
        </button>
        <div class="tab-spacer"></div>
        <button @click="showCreateDialog = true" class="btn btn-primary">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="12" y1="5" x2="12" y2="19"/>
            <line x1="5" y1="12" x2="19" y2="12"/>
          </svg>
          创建项目
        </button>
      </div>

      <!-- Search and Filter Bar (only for "全部项目" tab) -->
      <div v-if="activeTab === 'all'" class="filter-bar animate-in">
        <div class="search-input-wrapper">
          <svg class="search-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="11" cy="11" r="8"/>
            <line x1="21" y1="21" x2="16.65" y2="16.65"/>
          </svg>
          <input
            v-model="filters.keyword"
            type="text"
            placeholder="搜索项目名称..."
            class="search-input"
            @input="debouncedSearch"
          />
        </div>
        <select v-model="filters.isArchived" class="filter-select" @change="loadProjects">
          <option :value="undefined">全部状态</option>
          <option :value="false">活动中</option>
          <option :value="true">已归档</option>
        </select>
        <input
          v-model="filters.ownerName"
          type="text"
          placeholder="所有者名称..."
          class="filter-input"
          @input="debouncedSearch"
        />
        <button v-if="hasActiveFilters" class="btn btn-text" @click="clearFilters">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="18" y1="6" x2="6" y2="18"/>
            <line x1="6" y1="6" x2="18" y2="18"/>
          </svg>
          清除
        </button>
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="loading-state animate-in">
        <div class="loading-spinner"></div>
        <span>加载中...</span>
      </div>

      <!-- Error State -->
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

      <!-- Empty State -->
      <div v-else-if="projects.length === 0" class="empty-state animate-in">
        <div class="empty-illustration">
          <svg width="80" height="80" viewBox="0 0 80 80" fill="none">
            <rect x="10" y="20" width="60" height="50" rx="8" stroke="currentColor" stroke-width="2" stroke-dasharray="4 4"/>
            <rect x="20" y="10" width="40" height="30" rx="4" fill="var(--color-bg-muted)" stroke="currentColor" stroke-width="2"/>
            <line x1="30" y1="20" x2="50" y2="20" stroke="currentColor" stroke-width="2"/>
            <line x1="30" y1="28" x2="45" y2="28" stroke="currentColor" stroke-width="2"/>
          </svg>
        </div>
        <h3>暂无项目</h3>
        <p v-if="activeTab === 'my'">您还没有参与任何项目</p>
        <p v-else-if="hasActiveFilters">没有符合筛选条件的项目</p>
        <p v-else>创建您的第一个项目开始吧</p>
        <button @click="showCreateDialog = true" class="btn btn-primary">
          创建第一个项目
        </button>
      </div>

      <!-- Projects Grid -->
      <div v-else>
        <div class="projects-grid">
          <ProjectCard
            v-for="(project, index) in projects"
            :key="project.id"
            :project="project"
            :style="{ animationDelay: `${index * 0.05}s` }"
            class="animate-in"
            @click="(project) => goToProject(project.id)"
          />
        </div>

        <!-- Pagination -->
        <div v-if="totalPages > 1" class="pagination">
          <button
            class="page-btn"
            :disabled="currentPage <= 1"
            @click="goToPage(currentPage - 1)"
          >
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="15 18 9 12 15 6"/>
            </svg>
          </button>
          <div class="page-info">
            <span class="page-current">{{ currentPage }}</span>
            <span class="page-separator">/</span>
            <span class="page-total">{{ totalPages }}</span>
          </div>
          <button
            class="page-btn"
            :disabled="currentPage >= totalPages"
            @click="goToPage(currentPage + 1)"
          >
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="9 18 15 12 9 6"/>
            </svg>
          </button>
        </div>
      </div>
    </main>

    <!-- Create Project Dialog -->
    <CreateProjectDialog
      v-model:visible="showCreateDialog"
      @success="onProjectCreated"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { getProjects } from '../../api/project'
import CreateProjectDialog from '../../components/project/CreateProjectDialog.vue'
import ProjectCard from '../../components/project/ProjectCard.vue'
import AppHeader from '../../components/AppHeader.vue'
import type { Project, ApiResponse, PaginatedProjects } from '../../api/project'

const router = useRouter()
const route = useRoute()
const projects = ref<Project[]>([])
const loading = ref(false)
const error = ref<string | null>(null)
const showCreateDialog = ref(false)
const activeTab = ref<'my' | 'all'>('my')

// Pagination state
const currentPage = ref(1)
const pageSize = ref(20)
const totalPages = ref(1)

// Filter state
const filters = ref({
  keyword: '',
  ownerName: '',
  isArchived: undefined as boolean | undefined
})

// Debounce timer
let searchTimer: number | null = null

const hasActiveFilters = computed(() => {
  return filters.value.keyword !== '' ||
         filters.value.ownerName !== '' ||
         filters.value.isArchived !== undefined
})

onMounted(async () => {
  if (route.name === 'ProjectNew' || route.query.new === 'true') {
    showCreateDialog.value = true
  }
  await loadProjects()
})

async function loadProjects() {
  loading.value = true
  error.value = null
  try {
    const params: any = {
      page: currentPage.value,
      page_size: pageSize.value
    }
    if (activeTab.value === 'my') {
      params.my = 'true'
    }
    if (filters.value.keyword) {
      params.keyword = filters.value.keyword
    }
    if (filters.value.ownerName) {
      params.owner_name = filters.value.ownerName
    }
    if (filters.value.isArchived !== undefined) {
      params.is_archived = filters.value.isArchived
    }
    const res = await getProjects(params) as ApiResponse<PaginatedProjects>
    if (res.code === 0) {
      projects.value = res.data.projects
      totalPages.value = res.data.total_pages
    } else {
      error.value = res.message || '加载失败'
    }
  } catch (e: any) {
    error.value = e.message || '加载失败'
  } finally {
    loading.value = false
  }
}

function debouncedSearch() {
  if (searchTimer) {
    clearTimeout(searchTimer)
  }
  searchTimer = window.setTimeout(() => {
    currentPage.value = 1
    loadProjects()
  }, 300)
}

function switchTab(tab: 'my' | 'all') {
  activeTab.value = tab
  currentPage.value = 1
  filters.value = { keyword: '', ownerName: '', isArchived: undefined }
  loadProjects()
}

function goToPage(page: number) {
  if (page < 1 || page > totalPages.value) return
  currentPage.value = page
  loadProjects()
}

function clearFilters() {
  filters.value = { keyword: '', ownerName: '', isArchived: undefined }
  currentPage.value = 1
  loadProjects()
}

function goToProject(projectId: number) {
  router.push(`/projects/${projectId}`)
}

function onProjectCreated(newProject: Project) {
  projects.value.unshift(newProject)
  showCreateDialog.value = false
}
</script>

<style scoped>
.project-list-page {
  min-height: 100vh;
  background: var(--color-bg-base);
}

.page-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: var(--space-6);
}

/* Tab Bar */
.tab-bar {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  background: var(--color-bg-surface);
  padding: var(--space-2);
  border-radius: var(--radius-xl);
  border: 1px solid var(--color-border);
  margin-bottom: var(--space-4);
}

.tab-btn {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-3) var(--space-4);
  background: transparent;
  border: none;
  border-radius: var(--radius-lg);
  cursor: pointer;
  font-family: var(--font-body);
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  color: var(--color-text-secondary);
  transition: all var(--transition-base);
}

.tab-btn:hover {
  color: var(--color-text-primary);
  background: var(--color-bg-muted);
}

.tab-btn.active {
  background: var(--color-primary);
  color: white;
}

.tab-btn.active svg {
  opacity: 1;
}

.tab-btn svg {
  opacity: 0.7;
}

.tab-spacer {
  flex: 1;
}

/* Filter Bar */
.filter-bar {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-3);
  background: var(--color-bg-surface);
  border-radius: var(--radius-lg);
  border: 1px solid var(--color-border);
  margin-bottom: var(--space-4);
}

.search-input-wrapper {
  position: relative;
  flex: 1;
}

.search-icon {
  position: absolute;
  left: var(--space-3);
  top: 50%;
  transform: translateY(-50%);
  color: var(--color-text-muted);
}

.search-input {
  width: 100%;
  padding: var(--space-2) var(--space-3);
  padding-left: calc(var(--space-3) + 20px);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  font-family: var(--font-body);
  font-size: var(--font-size-sm);
  background: var(--color-bg-base);
  transition: all var(--transition-base);
}

.search-input:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 2px var(--color-primary-light);
}

.filter-select,
.filter-input {
  padding: var(--space-2) var(--space-3);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  font-family: var(--font-body);
  font-size: var(--font-size-sm);
  background: var(--color-bg-base);
  transition: all var(--transition-base);
  min-width: 120px;
}

.filter-input {
  min-width: 140px;
}

.filter-select:focus,
.filter-input:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 2px var(--color-primary-light);
}

/* Buttons */
.btn {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-4);
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

.btn-primary:hover {
  background: var(--color-primary-hover);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(91, 95, 199, 0.4);
}

.btn-text {
  background: transparent;
  color: var(--color-text-muted);
  padding: var(--space-2);
}

.btn-text:hover {
  color: var(--color-text-primary);
  background: var(--color-bg-muted);
}

/* Loading State */
.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--space-16);
  color: var(--color-text-muted);
  gap: var(--space-4);
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

/* Error State */
.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--space-16);
  color: var(--color-danger);
  gap: var(--space-4);
  text-align: center;
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

/* Empty State */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--space-16);
  text-align: center;
}

.empty-illustration {
  color: var(--color-text-muted);
  margin-bottom: var(--space-6);
}

.empty-state h3 {
  font-family: var(--font-display);
  font-size: var(--font-size-xl);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-primary);
  margin: 0 0 var(--space-2);
}

.empty-state p {
  color: var(--color-text-muted);
  margin: 0 0 var(--space-6);
}

/* Projects Grid */
.projects-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: var(--space-5);
}

/* Project Card */
.project-card {
  position: relative;
  background: var(--color-bg-surface);
  border-radius: var(--radius-xl);
  padding: var(--space-6);
  box-shadow: var(--shadow-card);
  cursor: pointer;
  transition: all var(--transition-base);
  border: 1px solid var(--color-border);
  overflow: hidden;
}

.project-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-card-hover);
  border-color: var(--color-primary-light);
}

.project-card:hover .card-accent {
  opacity: 1;
}

/* Card Accent */
.card-accent {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: linear-gradient(90deg, var(--color-primary) 0%, var(--color-primary-light) 100%);
  opacity: 0;
  transition: opacity var(--transition-base);
}

/* Card Header */
.project-card-header {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  margin-bottom: var(--space-3);
}

.project-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  background: var(--color-bg-muted);
  border-radius: var(--radius-lg);
  color: var(--color-primary);
}

.project-name {
  flex: 1;
  font-family: var(--font-display);
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-primary);
  margin: 0;
}

.archived-badge {
  padding: var(--space-1) var(--space-3);
  background: var(--color-bg-muted);
  color: var(--color-text-muted);
  border-radius: var(--radius-full);
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-medium);
}

/* Project Description */
.project-description {
  color: var(--color-text-secondary);
  font-size: var(--font-size-sm);
  line-height: var(--line-height-relaxed);
  margin: 0 0 var(--space-3);
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* Card Footer */
.project-card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: var(--space-4);
  border-top: 1px solid var(--color-border);
}

.created-at {
  font-size: var(--font-size-xs);
  color: var(--color-text-muted);
}

/* Pagination */
.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-3);
  margin-top: var(--space-6);
}

.page-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  background: var(--color-bg-surface);
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: all var(--transition-base);
}

.page-btn:hover:not(:disabled) {
  border-color: var(--color-primary);
  color: var(--color-primary);
}

.page-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.page-info {
  display: flex;
  align-items: center;
  gap: var(--space-1);
  font-size: var(--font-size-sm);
}

.page-current {
  font-weight: var(--font-weight-semibold);
  color: var(--color-primary);
}

.page-separator {
  color: var(--color-text-muted);
}

.page-total {
  color: var(--color-text-muted);
}

/* Responsive */
@media (max-width: 768px) {
  .page-content {
    padding: var(--space-4);
  }

  .tab-bar {
    flex-wrap: wrap;
  }

  .tab-spacer {
    width: 100%;
    height: 0;
  }

  .btn {
    width: 100%;
    justify-content: center;
  }

  .filter-bar {
    flex-wrap: wrap;
  }

  .search-input-wrapper {
    width: 100%;
  }

  .filter-select,
  .filter-input {
    flex: 1;
    min-width: 0;
  }

  .projects-grid {
    grid-template-columns: 1fr;
  }
}

/* Large screens - minimum 4 columns */
@media (min-width: 1200px) {
  .projects-grid {
    grid-template-columns: repeat(4, 1fr);
  }
}
</style>
