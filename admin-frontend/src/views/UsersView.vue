<template>
  <div class="users-view">
    <el-card class="box-card" shadow="hover">
      <template #header>
        <div class="card-header">
          <h1 class="page-title">用户管理</h1>
          <div class="header-actions">
            <el-input
              v-model="searchKeyword"
              placeholder="搜索用户名或邮箱"
              prefix-icon="Search"
              clearable
              style="width: 240px"
              @keyup.enter="handleSearch"
              @clear="handleSearch"
            />
            <el-button type="primary" icon="Plus" @click="openAddDialog" style="margin-left: 12px">
              添加用户
            </el-button>
          </div>
        </div>
      </template>

      <el-table
        v-loading="loading"
        :data="users"
        style="width: 100%"
        border
        stripe
      >
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="username" label="用户名" min-width="120" />
        <el-table-column prop="email" label="邮箱" min-width="180" />
        <el-table-column prop="role" label="角色" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getRoleTagType(row.role)">
              {{ getRoleDisplayName(row.role) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'danger'">
              {{ row.status === 'active' ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="260" align="center" fixed="right">
          <template #default="{ row }">
            <div class="action-buttons">
              <el-button
                type="primary"
                link
                size="small"
                @click="openEditDialog(row)"
              >
                编辑
              </el-button>
              <el-button
                type="warning"
                link
                size="small"
                @click="openResetPasswordDialog(row)"
              >
                重置密码
              </el-button>
              <el-button
                type="danger"
                link
                size="small"
                @click="handleDelete(row)"
              >
                删除
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 创建/编辑用户对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑用户' : '添加用户'"
      width="500px"
      @closed="resetForm"
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="用户名" prop="username">
          <el-input v-model="formData.username" placeholder="请输入用户名" :disabled="isEdit" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="formData.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item v-if="!isEdit" label="密码" prop="password">
          <el-input
            v-model="formData.password"
            type="password"
            placeholder="请输入密码（至少6位）"
            show-password
          />
        </el-form-item>
        <el-form-item label="角色" prop="role">
          <el-select v-model="formData.role" placeholder="请选择角色" style="width: 100%">
            <el-option label="管理员" value="admin" />
            <el-option label="成员" value="member" />
            <el-option label="查看者" value="viewer" />
          </el-select>
        </el-form-item>
        <el-form-item v-if="isEdit" label="状态" prop="status">
          <el-select v-model="formData.status" placeholder="请选择状态" style="width: 100%">
            <el-option label="启用" value="active" />
            <el-option label="禁用" value="disabled" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit" :loading="submitting">
            确定
          </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 重置密码对话框 -->
    <el-dialog
      v-model="passwordDialogVisible"
      title="重置密码"
      width="400px"
    >
      <el-form
        ref="passwordFormRef"
        :model="passwordFormData"
        :rules="passwordRules"
        label-width="100px"
      >
        <el-form-item label="用户">
          <span>{{ editingUser?.username }}</span>
        </el-form-item>
        <el-form-item label="新密码" prop="new_password">
          <el-input
            v-model="passwordFormData.new_password"
            type="password"
            placeholder="请输入新密码（至少6位）"
            show-password
          />
        </el-form-item>
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input
            v-model="passwordFormData.confirmPassword"
            type="password"
            placeholder="请再次输入新密码"
            show-password
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="passwordDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleResetPassword" :loading="resettingPassword">
            确定
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getUsers, createUser, updateUser, deleteUser, resetUserPassword } from '@/services/user'
import type { User, UpdateUserRequest } from '@/types/api'
import type { FormInstance, FormRules } from 'element-plus'

const loading = ref(false)
const users = ref<User[]>([])
const dialogVisible = ref(false)
const passwordDialogVisible = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const resettingPassword = ref(false)
const formRef = ref<FormInstance>()
const passwordFormRef = ref<FormInstance>()

// 分页
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const searchKeyword = ref('')

// 编辑状态
const editingId = ref<number | null>(null)
const editingUser = ref<User | null>(null)

// 用户表单数据
interface UserFormData {
  username: string
  email: string
  password: string
  role: 'admin' | 'member' | 'viewer'
  status: 'active' | 'disabled'
}

const formData = reactive<UserFormData>({
  username: '',
  email: '',
  password: '',
  role: 'member',
  status: 'active',
})

const passwordFormData = reactive({
  new_password: '',
  confirmPassword: '',
})

const validateConfirmPassword = (_rule: unknown, value: string, callback: (error?: Error) => void) => {
  if (value !== passwordFormData.new_password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const rules = reactive<FormRules>({
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 50, message: '用户名长度在 3-50 个字符', trigger: 'blur' },
    { pattern: /^[a-zA-Z0-9_]+$/, message: '用户名只能包含字母、数字和下划线', trigger: 'blur' },
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入有效的邮箱地址', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 100, message: '密码长度至少 6 个字符', trigger: 'blur' },
  ],
  role: [
    { required: true, message: '请选择角色', trigger: 'change' },
  ],
  status: [
    { required: true, message: '请选择状态', trigger: 'change' },
  ],
})

const passwordRules = reactive<FormRules>({
  new_password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, max: 100, message: '密码长度至少 6 个字符', trigger: 'blur' },
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' },
  ],
})

onMounted(() => {
  fetchUsers()
})

const fetchUsers = async () => {
  loading.value = true
  try {
    const response = await getUsers({
      page: currentPage.value,
      page_size: pageSize.value,
      keyword: searchKeyword.value,
    })
    users.value = response.data
    total.value = response.meta.total_count
  } catch (error) {
    const message = error instanceof Error ? error.message : '获取用户列表失败'
    ElMessage.error(message)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  currentPage.value = 1
  fetchUsers()
}

const handleSizeChange = () => {
  currentPage.value = 1
  fetchUsers()
}

const handleCurrentChange = () => {
  fetchUsers()
}

const formatDate = (dateString: string) => {
  if (!dateString) return '-'
  return new Date(dateString).toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}

const getRoleDisplayName = (role: string) => {
  const roleMap: Record<string, string> = {
    admin: '管理员',
    member: '成员',
    viewer: '查看者',
  }
  return roleMap[role] || role
}

const getRoleTagType = (role: string) => {
  const typeMap: Record<string, string> = {
    admin: 'danger',
    member: 'primary',
    viewer: 'info',
  }
  return typeMap[role] || 'info'
}

const resetForm = () => {
  if (formRef.value) {
    formRef.value.resetFields()
  }
  formData.username = ''
  formData.email = ''
  formData.password = ''
  formData.role = 'member'
  formData.status = 'active'
  editingId.value = null
  isEdit.value = false
}

const openAddDialog = () => {
  isEdit.value = false
  dialogVisible.value = true
}

const openEditDialog = (row: User) => {
  isEdit.value = true
  editingId.value = row.id
  formData.username = row.username || ''
  formData.email = row.email || ''
  formData.role = (row.role as 'admin' | 'member' | 'viewer') || 'member'
  formData.status = (row.status as 'active' | 'disabled') || 'active'
  dialogVisible.value = true
}

const openResetPasswordDialog = (row: User) => {
  editingUser.value = row
  passwordFormData.new_password = ''
  passwordFormData.confirmPassword = ''
  passwordDialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        if (isEdit.value && editingId.value !== null) {
          const updateData: UpdateUserRequest = {
            email: formData.email,
            role: formData.role,
            status: formData.status,
          }
          await updateUser(editingId.value as number, updateData)
          ElMessage.success('更新成功')
        } else {
          await createUser({
            username: formData.username,
            email: formData.email,
            password: formData.password,
            role: formData.role,
          })
          ElMessage.success('创建成功')
        }
        dialogVisible.value = false
        fetchUsers()
      } catch (error) {
        const message = error instanceof Error ? error.message : '操作失败'
        ElMessage.error(message)
      } finally {
        submitting.value = false
      }
    }
  })
}

const handleResetPassword = async () => {
  if (!passwordFormRef.value || editingId.value === null) return

  await passwordFormRef.value.validate(async (valid) => {
    if (valid) {
      resettingPassword.value = true
      try {
        await resetUserPassword(editingId.value as number, {
          new_password: passwordFormData.new_password,
        })
        ElMessage.success('密码重置成功')
        passwordDialogVisible.value = false
      } catch (error) {
        const message = error instanceof Error ? error.message : '重置密码失败'
        ElMessage.error(message)
      } finally {
        resettingPassword.value = false
      }
    }
  })
}

const handleDelete = (row: User) => {
  ElMessageBox.confirm(
    `确定要删除用户 "${row.username}" 吗？此操作不可恢复。`,
    '删除确认',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(async () => {
    try {
      await deleteUser(row.id)
      ElMessage.success('删除成功')
      fetchUsers()
    } catch (error) {
      const message = error instanceof Error ? error.message : '删除失败'
      ElMessage.error(message)
    }
  }).catch(() => {})
}
</script>

<style scoped>
.users-view {
  background-color: #f3f4f6;
}

.box-card {
  border-radius: 8px;
  border: none;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
}

.header-actions {
  display: flex;
  align-items: center;
}

.page-title {
  font-size: 1.5rem;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}

.action-buttons {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  flex-wrap: nowrap;
}
</style>
