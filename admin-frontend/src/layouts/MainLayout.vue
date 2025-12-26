<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import {
  Fold,
  Expand,
  User,
  Odometer,
  FolderOpened,
  ChatDotRound,
  Document,
  Setting,
} from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

// 侧边栏折叠状态
const isCollapse = ref(false)

// 菜单项配置
const menuItems = [
  {
    index: '/dashboard',
    title: '仪表板',
    icon: Odometer,
  },
  {
    index: '/projects',
    title: '项目管理',
    icon: FolderOpened,
  },
  {
    index: '/languages',
    title: '语言管理',
    icon: ChatDotRound,
  },
  {
    index: '/translations',
    title: '翻译管理',
    icon: Document,
  },
  {
    index: '/settings',
    title: '系统设置',
    icon: Setting,
  },
]

// 当前激活的菜单项
const activeMenu = computed(() => {
  return route.path
})

// 页面标题
const pageTitle = computed(() => {
  return route.meta.title as string || '仪表板'
})

// 切换侧边栏
const toggleSidebar = () => {
  isCollapse.value = !isCollapse.value
}

// 处理菜单点击
const handleMenuSelect = (index: string) => {
  router.push(index)
}

// 退出登录
const handleLogout = () => {
  authStore.logout()
}
</script>

<template>
  <el-container class="layout-container">
    <!-- 侧边栏 -->
    <el-aside :width="isCollapse ? '64px' : '200px'" class="layout-aside">
      <div class="logo-container">
        <div v-if="!isCollapse" class="logo-text">
          <h2>i18n-flow</h2>
        </div>
        <div v-else class="logo-icon">
          <span>i18n</span>
        </div>
      </div>

      <el-menu
        :default-active="activeMenu"
        :collapse="isCollapse"
        :collapse-transition="false"
        class="layout-menu"
        background-color="#001529"
        text-color="#ffffff"
        active-text-color="#1890ff"
        @select="handleMenuSelect"
      >
        <el-menu-item
          v-for="item in menuItems"
          :key="item.index"
          :index="item.index"
        >
          <el-icon><component :is="item.icon" /></el-icon>
          <template #title>{{ item.title }}</template>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <!-- 右侧容器 -->
    <el-container>
      <!-- 头部 -->
      <el-header class="layout-header">
        <div class="header-left">
          <el-button
            :icon="isCollapse ? Expand : Fold"
            text
            @click="toggleSidebar"
          />
          <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/dashboard' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item>{{ pageTitle }}</el-breadcrumb-item>
          </el-breadcrumb>
        </div>

        <div class="header-right">
          <el-dropdown @command="handleLogout">
            <div class="user-info">
              <el-avatar :size="32" :src="''">
                <el-icon><User /></el-icon>
              </el-avatar>
              <span class="username">{{ authStore.user?.username }}</span>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item disabled>
                  <div class="user-detail">
                    <p><strong>用户 ID:</strong> {{ authStore.user?.id }}</p>
                    <p><strong>用户名:</strong> {{ authStore.user?.username }}</p>
                  </div>
                </el-dropdown-item>
                <el-dropdown-item divided command="logout">
                  退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <!-- 主内容区 -->
      <el-main class="layout-main">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<style scoped>
.layout-container {
  height: 100vh;
}

.layout-aside {
  background-color: #001529;
  transition: width 0.3s;
  overflow-x: hidden;
}

.logo-container {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.logo-text h2 {
  margin: 0;
  color: #ffffff;
  font-size: 20px;
  font-weight: 600;
  letter-spacing: 1px;
}

.logo-icon {
  color: #1890ff;
  font-size: 16px;
  font-weight: 700;
}

.layout-menu {
  border-right: none;
  height: calc(100vh - 60px);
}

.layout-menu:not(.el-menu--collapse) {
  width: 200px;
}

.layout-header {
  background: #ffffff;
  border-bottom: 1px solid #f0f0f0;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.header-right {
  display: flex;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 4px;
  transition: background-color 0.3s;
}

.user-info:hover {
  background-color: #f5f5f5;
}

.username {
  font-size: 14px;
  color: #303133;
}

.user-detail {
  padding: 8px 0;
}

.user-detail p {
  margin: 4px 0;
  font-size: 13px;
  color: #606266;
  white-space: nowrap;
}

.layout-main {
  background-color: #f5f7fa;
  padding: 24px;
  overflow-y: auto;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .layout-main {
    padding: 16px;
  }

  .username {
    display: none;
  }
}
</style>
