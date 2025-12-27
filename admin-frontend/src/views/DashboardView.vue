<script setup lang="ts">
import { computed } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import { getDashboardStats } from '@/services/dashboardService'
import {
  Folder,
  Document,
  MessageBox,
  ChatDotRound,
} from '@element-plus/icons-vue'

// 使用 vue-query 获取仪表板统计数据
const {
  data: stats,
  isLoading,
  isError,
  error,
} = useQuery({
  queryKey: ['dashboardStats'],
  queryFn: getDashboardStats,
})

// 统计卡片配置
const statCards = computed(() => [
  {
    title: '项目总数',
    value: stats.value?.TotalProjects ?? 0,
    icon: Folder,
    color: '#409eff',
    bgColor: '#ecf5ff',
  },
  {
    title: '语言总数',
    value: stats.value?.TotalLanguages ?? 0,
    icon: ChatDotRound,
    color: '#67c23a',
    bgColor: '#f0f9ff',
  },
  {
    title: '翻译键总数',
    value: stats.value?.TotalKeys ?? 0,
    icon: Document,
    color: '#e6a23c',
    bgColor: '#fdf6ec',
  },
  {
    title: '翻译总数',
    value: stats.value?.TotalTranslations ?? 0,
    icon: MessageBox,
    color: '#f56c6c',
    bgColor: '#fef0f0',
  },
])
</script>

<template>
  <div class="dashboard-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1>仪表板</h1>
      <p class="subtitle">系统概览和统计信息</p>
    </div>

    <!-- 加载中状态 -->
    <el-row v-if="isLoading" :gutter="20" class="stats-grid">
      <el-col v-for="i in 4" :key="i" :xs="24" :sm="12" :lg="6">
        <el-card class="stat-card">
          <el-skeleton animated>
            <template #template>
              <el-skeleton-item variant="circle" style="width: 48px; height: 48px" />
              <el-skeleton-item variant="text" style="width: 60%; margin-top: 16px" />
              <el-skeleton-item variant="h1" style="width: 40%; margin-top: 8px" />
            </template>
          </el-skeleton>
        </el-card>
      </el-col>
    </el-row>

    <!-- 错误状态 -->
    <el-alert
      v-else-if="isError"
      type="error"
      :title="error?.message || '加载统计数据失败'"
      :description="'请检查网络连接或稍后重试'"
      show-icon
      :closable="false"
    />

    <!-- 统计卡片 -->
    <el-row v-else :gutter="20" class="stats-grid">
      <el-col
        v-for="card in statCards"
        :key="card.title"
        :xs="24"
        :sm="12"
        :lg="6"
      >
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <div class="stat-icon" :style="{ backgroundColor: card.bgColor }">
              <el-icon :size="32" :color="card.color">
                <component :is="card.icon" />
              </el-icon>
            </div>
            <div class="stat-info">
              <p class="stat-title">{{ card.title }}</p>
              <h2 class="stat-value" :style="{ color: card.color }">
                {{ card.value.toLocaleString() }}
              </h2>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<style scoped>
.dashboard-page {
  max-width: 1400px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 24px;
}

.page-header h1 {
  margin: 0;
  font-size: 28px;
  color: #303133;
  font-weight: 600;
}

.subtitle {
  margin: 8px 0 0;
  font-size: 14px;
  color: #909399;
}

.stats-grid {
  margin-bottom: 24px;
}

.stat-card {
  border-radius: 8px;
  transition: all 0.3s;
  cursor: default;
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 64px;
  height: 64px;
  border-radius: 12px;
  flex-shrink: 0;
}

.stat-info {
  flex: 1;
}

.stat-title {
  margin: 0 0 8px;
  font-size: 14px;
  color: #909399;
  font-weight: 500;
}

.stat-value {
  margin: 0;
  font-size: 32px;
  font-weight: 700;
  line-height: 1;
}

/* 响应式调整 */
@media (max-width: 768px) {
  .page-header h1 {
    font-size: 24px;
  }

  .stat-value {
    font-size: 28px;
  }
}
</style>
