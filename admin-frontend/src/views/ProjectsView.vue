<script setup lang="ts">
import { ref, computed } from 'vue'
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import {
  getProjects,
  createProject,
  updateProject,
  deleteProject,
} from '@/services/projectService'
import type {
  Project,
  CreateProjectRequest,
  UpdateProjectRequest,
} from '@/types/api'
import {
  Plus,
  Edit,
  Delete,
  Search,
  Refresh,
} from '@element-plus/icons-vue'

// ============ 状态管理 ============
const queryClient = useQueryClient()

// 搜索和分页状态
const searchKeyword = ref('')
const currentPage = ref(1)
const pageSize = ref(10)

// 对话框状态
const dialogVisible = ref(false)
const dialogTitle = ref('创建项目')
const isEditMode = ref(false)
const currentEditId = ref<number | null>(null)

// 表单引用和数据
const formRef = ref<FormInstance>()
const formData = ref<CreateProjectRequest | UpdateProjectRequest>({
  name: '',
  description: '',
})

// 表单验证规则
const formRules: FormRules = {
  name: [
    { required: true, message: '请输入项目名称', trigger: 'blur' },
    { min: 2, max: 100, message: '项目名称长度在 2 到 100 个字符', trigger: 'blur' },
  ],
}

// ============ 数据获取 ============
const queryParams = computed(() => ({
  page: currentPage.value,
  page_size: pageSize.value,
  keyword: searchKeyword.value || undefined,
}))

const {
  data: projectsData,
  isLoading,
  isError,
  error,
  refetch,
} = useQuery({
  queryKey: ['projects', queryParams],
  queryFn: () => getProjects(queryParams.value),
})

const projects = computed(() => projectsData.value?.data || [])
const totalCount = computed(() => projectsData.value?.meta?.total_count || 0)

// ============ CRUD 操作 ============

// 创建项目
const createMutation = useMutation({
  mutationFn: createProject,
  onSuccess: () => {
    ElMessage.success('项目创建成功')
    dialogVisible.value = false
    queryClient.invalidateQueries({ queryKey: ['projects'] })
  },
  onError: (err: Error) => {
    ElMessage.error(err.message || '创建项目失败')
  },
})

// 更新项目
const updateMutation = useMutation({
  mutationFn: ({ id, data }: { id: number; data: UpdateProjectRequest }) =>
    updateProject(id, data),
  onSuccess: () => {
    ElMessage.success('项目更新成功')
    dialogVisible.value = false
    queryClient.invalidateQueries({ queryKey: ['projects'] })
  },
  onError: (err: Error) => {
    ElMessage.error(err.message || '更新项目失败')
  },
})

// 删除项目
const deleteMutation = useMutation({
  mutationFn: deleteProject,
  onSuccess: () => {
    ElMessage.success('项目删除成功')
    queryClient.invalidateQueries({ queryKey: ['projects'] })
  },
  onError: (err: Error) => {
    ElMessage.error(err.message || '删除项目失败')
  },
})

// ============ 事件处理 ============

// 搜索
const handleSearch = () => {
  currentPage.value = 1 // 搜索时重置到第一页
}

// 清空搜索
const handleClearSearch = () => {
  searchKeyword.value = ''
  handleSearch()
}

// 刷新列表
const handleRefresh = () => {
  refetch()
}

// 分页变化
const handlePageChange = (page: number) => {
  currentPage.value = page
}

const handlePageSizeChange = (size: number) => {
  pageSize.value = size
  currentPage.value = 1
}

// 打开创建对话框
const handleCreate = () => {
  isEditMode.value = false
  dialogTitle.value = '创建项目'
  formData.value = {
    name: '',
    description: '',
  }
  dialogVisible.value = true
}

// 打开编辑对话框
const handleEdit = (project: Project) => {
  isEditMode.value = true
  dialogTitle.value = '编辑项目'
  currentEditId.value = project.id
  formData.value = {
    name: project.name,
    description: project.description || '',
    status: project.status,
  }
  dialogVisible.value = true
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate((valid) => {
    if (valid) {
      if (isEditMode.value && currentEditId.value) {
        updateMutation.mutate({
          id: currentEditId.value,
          data: formData.value as UpdateProjectRequest,
        })
      } else {
        createMutation.mutate(formData.value as CreateProjectRequest)
      }
    }
  })
}

// 取消表单
const handleCancel = () => {
  dialogVisible.value = false
  formRef.value?.resetFields()
}

// 删除项目
const handleDelete = (project: Project) => {
  ElMessageBox.confirm(
    `确定要删除项目 "${project.name}" 吗？此操作不可恢复。`,
    '删除确认',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  )
    .then(() => {
      deleteMutation.mutate(project.id)
    })
    .catch(() => {
      // 用户取消删除
    })
}

// 切换项目状态
const handleToggleStatus = (project: Project) => {
  const newStatus = project.status === 'active' ? 'archived' : 'active'
  updateMutation.mutate({
    id: project.id,
    data: { status: newStatus },
  })
}

// 格式化日期
const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}
</script>

<template>
  <div class="projects-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div>
        <h1>项目管理</h1>
        <p class="subtitle">管理所有翻译项目</p>
      </div>
      <el-button type="primary" :icon="Plus" @click="handleCreate">
        创建项目
      </el-button>
    </div>

    <!-- 搜索和操作栏 -->
    <el-card class="search-card" shadow="never">
      <div class="search-bar">
        <div class="search-input-group">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索项目名称或描述"
            :prefix-icon="Search"
            clearable
            @clear="handleClearSearch"
            @keyup.enter="handleSearch"
          />
          <el-button type="primary" :icon="Search" @click="handleSearch">
            搜索
          </el-button>
        </div>
        <el-button :icon="Refresh" @click="handleRefresh">刷新</el-button>
      </div>
    </el-card>

    <!-- 项目列表 -->
    <el-card class="table-card" shadow="never">
      <!-- 加载中状态 -->
      <div v-if="isLoading" class="loading-container">
        <el-skeleton :rows="5" animated />
      </div>

      <!-- 错误状态 -->
      <el-alert
        v-else-if="isError"
        type="error"
        :title="error?.message || '加载项目列表失败'"
        :description="'请检查网络连接或稍后重试'"
        show-icon
        :closable="false"
      />

      <!-- 数据表格 -->
      <div v-else>
        <el-table
          :data="projects"
          style="width: 100%"
          :empty-text="searchKeyword ? '没有找到匹配的项目' : '暂无项目数据'"
        >
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column prop="name" label="项目名称" min-width="150">
            <template #default="{ row }">
              <div class="project-name">
                <strong>{{ row.name }}</strong>
                <span class="project-slug">{{ row.slug }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column
            prop="description"
            label="描述"
            min-width="200"
            show-overflow-tooltip
          />
          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <el-tag
                :type="row.status === 'active' ? 'success' : 'info'"
                size="small"
                @click="handleToggleStatus(row)"
                style="cursor: pointer"
              >
                {{ row.status === 'active' ? '活跃' : '已归档' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="创建时间" width="160">
            <template #default="{ row }">
              {{ formatDate(row.created_at) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="180" fixed="right">
            <template #default="{ row }">
              <el-button
                type="primary"
                size="small"
                :icon="Edit"
                link
                @click="handleEdit(row)"
              >
                编辑
              </el-button>
              <el-button
                type="danger"
                size="small"
                :icon="Delete"
                link
                @click="handleDelete(row)"
              >
                删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>

        <!-- 分页 -->
        <div class="pagination-container">
          <el-pagination
            v-model:current-page="currentPage"
            v-model:page-size="pageSize"
            :page-sizes="[10, 20, 50, 100]"
            :total="totalCount"
            layout="total, sizes, prev, pager, next, jumper"
            @current-change="handlePageChange"
            @size-change="handlePageSizeChange"
          />
        </div>
      </div>
    </el-card>

    <!-- 创建/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="500px"
      :close-on-click-modal="false"
      @close="handleCancel"
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="80px"
      >
        <el-form-item label="项目名称" prop="name">
          <el-input
            v-model="formData.name"
            placeholder="请输入项目名称"
            clearable
          />
        </el-form-item>
        <el-form-item label="项目描述" prop="description">
          <el-input
            v-model="formData.description"
            type="textarea"
            :rows="4"
            placeholder="请输入项目描述（可选）"
            clearable
          />
        </el-form-item>
        <el-form-item v-if="isEditMode" label="项目状态" prop="status">
          <el-radio-group v-model="(formData as UpdateProjectRequest).status">
            <el-radio value="active">活跃</el-radio>
            <el-radio value="archived">已归档</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="handleCancel">取消</el-button>
        <el-button
          type="primary"
          :loading="createMutation.isPending.value || updateMutation.isPending.value"
          @click="handleSubmit"
        >
          {{ isEditMode ? '更新' : '创建' }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.projects-page {
  max-width: 1400px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
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

.search-card {
  margin-bottom: 16px;
}

.search-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
}

.search-input-group {
  display: flex;
  gap: 12px;
  flex: 1;
  max-width: 600px;
}

.table-card {
  border-radius: 8px;
}

.loading-container {
  padding: 20px;
}

.project-name {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.project-slug {
  font-size: 12px;
  color: #909399;
  font-family: monospace;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid #f0f0f0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }

  .page-header h1 {
    font-size: 24px;
  }

  .search-bar {
    flex-direction: column;
    align-items: stretch;
  }

  .search-input-group {
    max-width: 100%;
  }

  .pagination-container {
    justify-content: center;
  }

  :deep(.el-pagination) {
    flex-wrap: wrap;
    justify-content: center;
  }
}
</style>
