<template>
  <div
    class="project-card"
    @click="$emit('click', project)"
  >
    <!-- Card Accent -->
    <div class="card-accent"></div>

    <!-- Card Header -->
    <div class="project-card-header">
      <div class="project-icon">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <rect x="3" y="3" width="18" height="18" rx="2" ry="2"/>
          <line x1="3" y1="9" x2="21" y2="9"/>
          <line x1="9" y1="21" x2="9" y2="9"/>
        </svg>
      </div>
      <h3 class="project-name">{{ project.name }}</h3>
      <span v-if="project.is_archived" class="archived-badge">已归档</span>
    </div>

    <!-- Project Description -->
    <p class="project-description">{{ project.description || '暂无描述' }}</p>

    <!-- Card Footer -->
    <div class="project-card-footer">
      <div v-if="project.owner_name" class="project-meta">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
          <circle cx="12" cy="7" r="4"/>
        </svg>
        {{ project.owner_name }}
      </div>
      <span class="created-at">
        {{ formatDate(project.created_at) }}
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Project } from '../../api/project'

defineProps<{
  project: Project
}>()

defineEmits<{
  click: [project: Project]
}>()

function formatDate(dateStr: string): string {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN', { month: 'short', day: 'numeric' })
}
</script>

<style scoped>
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
  margin: 0 0 var(--space-4);
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

.project-meta {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
}

.project-meta svg {
  color: var(--color-primary);
}

.created-at {
  font-size: var(--font-size-xs);
  color: var(--color-text-muted);
}
</style>
