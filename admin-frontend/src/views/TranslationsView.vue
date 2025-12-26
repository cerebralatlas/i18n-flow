<template>
  <div class="translations-view">
    <div class="header">
      <h1>翻译管理</h1>
      <div class="header-actions">
        <select v-model="selectedProjectId" @change="handleProjectChange" class="project-selector">
          <option :value="null" disabled>选择项目</option>
          <option v-for="project in projects" :key="project.id" :value="project.id">
            {{ project.name }}
          </option>
        </select>
      </div>
    </div>

    <div v-if="selectedProjectId" class="content">
      <!-- Toolbar -->
      <div class="toolbar">
        <div class="search-box">
          <input
            v-model="searchKeyword"
            @input="handleSearch"
            type="text"
            placeholder="搜索翻译键..."
            class="search-input"
          />
        </div>
        <div class="toolbar-actions">
          <button @click="showAddKeyDialog = true" class="btn btn-primary">
            <span class="icon">+</span> 添加翻译键
          </button>
          <button @click="handleExport" class="btn btn-secondary" :disabled="loading">
            <span class="icon">↓</span> 导出
          </button>
          <button @click="showImportDialog = true" class="btn btn-secondary">
            <span class="icon">↑</span> 导入
          </button>
        </div>
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="loading">加载中...</div>

      <!-- Error State -->
      <div v-else-if="error" class="error">{{ error }}</div>

      <!-- Translation Matrix Table -->
      <div v-else-if="matrix && matrix.languages && matrix.rows" class="table-wrapper">
        <table class="translation-matrix">
          <thead>
            <tr>
              <th class="key-column">翻译键</th>
              <th class="context-column">上下文</th>
              <th v-for="lang in matrix.languages" :key="lang.id" class="language-column">
                {{ lang.name }}
                <span class="language-code">{{ lang.code }}</span>
              </th>
              <th class="actions-column">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="matrix.rows.length === 0">
              <td :colspan="(matrix.languages?.length || 0) + 3" class="empty-state">
                暂无翻译数据，点击"添加翻译键"开始
              </td>
            </tr>
            <tr v-for="(row, index) in matrix.rows" :key="row.key_name">
              <td class="key-column">
                <strong>{{ row.key_name }}</strong>
              </td>
              <td class="context-column">
                <span class="context-text">{{ row.context || '-' }}</span>
              </td>
              <td
                v-for="lang in matrix.languages"
                :key="lang.id"
                class="translation-cell"
                @click="editCell(row.key_name, lang)"
              >
                <div
                  v-if="
                    editingCell?.keyName === row.key_name && editingCell?.languageId === lang.id
                  "
                  class="cell-editing"
                >
                  <textarea
                    v-model="editingValue"
                    @blur="saveCell"
                    @keydown.enter.exact.prevent="saveCell"
                    @keydown.esc="cancelEdit"
                    ref="editInput"
                    class="cell-input"
                  ></textarea>
                </div>
                <div v-else class="cell-display">
                  <span v-if="row.translations[lang.code]?.value" class="cell-value">
                    {{ row.translations[lang.code].value }}
                  </span>
                  <span v-else class="cell-empty">未翻译</span>
                </div>
              </td>
              <td class="actions-column">
                <button
                  @click="handleDeleteKey(row.key_name)"
                  class="btn-icon btn-danger"
                  title="删除翻译键"
                >
                  ×
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div v-if="matrix && matrix.total_pages > 1" class="pagination">
        <button
          @click="changePage(matrix.page - 1)"
          :disabled="matrix.page <= 1"
          class="btn btn-secondary"
        >
          上一页
        </button>
        <span class="page-info">
          第 {{ matrix.page }} / {{ matrix.total_pages }} 页 (共 {{ matrix.total_count }} 条)
        </span>
        <button
          @click="changePage(matrix.page + 1)"
          :disabled="matrix.page >= matrix.total_pages"
          class="btn btn-secondary"
        >
          下一页
        </button>
      </div>
    </div>

    <div v-else class="empty-project">
      <p>请选择一个项目以管理翻译</p>
    </div>

    <!-- Add Key Dialog -->
    <div v-if="showAddKeyDialog" class="modal-overlay" @click.self="showAddKeyDialog = false">
      <div class="modal large-modal">
        <h2>添加翻译键</h2>
        <form @submit.prevent="handleAddKey">
          <div class="form-group">
            <label>翻译键名 <span class="required">*</span></label>
            <input v-model="newKey.keyName" type="text" required class="form-input" placeholder="例如: welcome.message" />
          </div>
          <div class="form-group">
            <label>上下文（可选）</label>
            <input v-model="newKey.context" type="text" class="form-input" placeholder="说明这个翻译键的使用场景" />
          </div>
          
          <div class="form-group">
            <label>翻译内容</label>
            <div class="language-inputs">
              <div v-for="lang in availableLanguages" :key="lang.id" class="language-input-row">
                <label class="language-label">
                  <span class="language-name">{{ lang.name }}</span>
                  <span class="language-code-badge">{{ lang.code }}</span>
                </label>
                <input 
                  v-model="newKey.translations[lang.code]" 
                  type="text" 
                  class="form-input"
                  :placeholder="`输入 ${lang.name} 翻译（可选）`"
                />
              </div>
            </div>
          </div>
          
          <div class="form-actions">
            <button type="button" @click="showAddKeyDialog = false" class="btn btn-secondary">
              取消
            </button>
            <button type="submit" class="btn btn-primary">添加</button>
          </div>
        </form>
      </div>
    </div>

    <!-- Import Dialog -->
    <div v-if="showImportDialog" class="modal-overlay" @click.self="showImportDialog = false">
      <div class="modal">
        <h2>导入翻译</h2>
        <div class="form-group">
          <label>选择JSON文件</label>
          <input type="file" @change="handleFileSelect" accept=".json" class="form-input" />
        </div>
        <div class="form-actions">
          <button @click="showImportDialog = false" class="btn btn-secondary">取消</button>
          <button @click="handleImport" :disabled="!importFile" class="btn btn-primary">
            导入
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick, watch, computed } from 'vue'
import { getTranslationMatrix, batchCreateTranslations, exportTranslations, importTranslations } from '@/services/translation'
import { getLanguages } from '@/services/language'
import type { TranslationMatrix, Language, BatchTranslationRequest, ImportTranslationsData } from '@/types/translation'
import type { Project } from '@/types/api'
import api from '@/services/api'

// State
const projects = ref<Project[]>([])
const selectedProjectId = ref<number | null>(null)
const matrix = ref<TranslationMatrix | null>(null)
const loading = ref(false)
const error = ref<string | null>(null)
const allLanguages = ref<Language[]>([])

// Search and Pagination
const searchKeyword = ref('')
const currentPage = ref(1)
const pageSize = ref(20)

// Editing
const editingCell = ref<{ keyName: string; languageId: number } | null>(null)
const editingValue = ref('')
const editInput = ref<HTMLTextAreaElement | null>(null)

// Dialogs
const showAddKeyDialog = ref(false)
const showImportDialog = ref(false)
const newKey = ref({ keyName: '', context: '', translations: {} as Record<string, string> })
const importFile = ref<File | null>(null)

// Available languages (active only)
const availableLanguages = computed(() => {
  return allLanguages.value.filter(l => l.status === 'active')
})

// Load projects on mount
onMounted(async () => {
  await loadProjects()
  await loadAllLanguages()
})

// Load all languages
const loadAllLanguages = async () => {
  try {
    allLanguages.value = await getLanguages()
  } catch (err: any) {
    console.error('Failed to load languages:', err)
  }
}

// Load projects list
const loadProjects = async () => {
  try {
    const response = await api.get('/projects', { params: { page: 1, page_size: 100 } })
    projects.value = response.data || []
  } catch (err: any) {
    console.error('Failed to load projects:', err)
  }
}

// Handle project change
const handleProjectChange = () => {
  currentPage.value = 1
  searchKeyword.value = ''
  loadMatrix()
}

// Load translation matrix
const loadMatrix = async () => {
  if (!selectedProjectId.value) return

  loading.value = true
  error.value = null

  try {
    // Get languages first
    const languages = await getLanguages()
    
    // Get matrix data from backend (returns map[string]map[string]string)
    const response = await api.get(`/translations/matrix/by-project/${selectedProjectId.value}`, {
      params: {
        page: currentPage.value,
        page_size: pageSize.value,
        keyword: searchKeyword.value || undefined
      }
    })
    
    console.log('Translation matrix API response:', response)
    
    // Extract data and meta from response
    const matrixData = response.data || {}
    const meta = response.meta || {}
    
    // Transform backend map structure to TranslationMatrix
    // Backend returns: { "key1": { "en": "value1", "zh": "value2" }, ...}
    const rows: any[] = []
    
    for (const [keyName, translations] of Object.entries(matrixData)) {
      if (keyName && typeof translations === 'object') {
        const translationCells: Record<string, any> = {}
        
        // Transform language code -> value to our cell structure
        for (const [langCode, value] of Object.entries(translations as Record<string, any>)) {
          const lang = languages.find(l => l.code === langCode)
          if (lang) {
            translationCells[langCode] = {
              language_id: lang.id,
              value: value || '',
              // We don't have ID from matrix endpoint, will need to fetch when editing
            }
          }
        }
        
        rows.push({
          key_name: keyName,
          context: '', // Backend doesn't return context in matrix view
          translations: translationCells
        })
      }
    }
    
    matrix.value = {
      languages: languages.filter(l => l.status === 'active'),
      rows: rows,
      total_count: meta.total_count || rows.length,
      page: meta.page || currentPage.value,
      page_size: meta.page_size || pageSize.value,
      total_pages: meta.total_pages || 1
    }
    
    console.log('Transformed matrix:', matrix.value)
    
  } catch (err: any) {
    error.value = err.message || '加载翻译数据失败'
    console.error('Failed to load translation matrix:', err)
    matrix.value = null
  } finally {
    loading.value = false
  }
}

// Search handler with debounce
let searchTimeout: number | null = null
const handleSearch = () => {
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = window.setTimeout(() => {
    currentPage.value = 1
    loadMatrix()
  }, 300)
}

// Pagination
const changePage = (page: number) => {
  currentPage.value = page
  loadMatrix()
}

// Edit cell
const editCell = (keyName: string, lang: Language) => {
  const row = matrix.value?.rows.find((r) => r.key_name === keyName)
  if (!row) return

  editingCell.value = { keyName, languageId: lang.id }
  editingValue.value = row.translations[lang.code]?.value || ''

  nextTick(() => {
    editInput.value?.[0]?.focus()
  })
}

// Save cell
const saveCell = async () => {
  if (!editingCell.value || !selectedProjectId.value) return

  const { keyName, languageId } = editingCell.value
  const row = matrix.value?.rows.find((r) => r.key_name === keyName)
  if (!row) return

  const lang = matrix.value?.languages.find((l) => l.id === languageId)
  if (!lang) return

  const existingCell = row.translations[lang.code]

  try {
    if (existingCell?.id) {
      // Update existing translation
      await api.put(`/translations/${existingCell.id}`, {
        project_id: selectedProjectId.value,
        language_id: languageId,
        key_name: keyName,
        value: editingValue.value,
        context: row.context,
      })
    } else {
      // Create new translation
      await api.post('/translations', {
        project_id: selectedProjectId.value,
        language_id: languageId,
        key_name: keyName,
        value: editingValue.value,
        context: row.context,
      })
    }

    // Refresh matrix
    await loadMatrix()
  } catch (err: any) {
    alert('保存失败: ' + (err.message || '未知错误'))
  } finally {
    cancelEdit()
  }
}

// Cancel edit
const cancelEdit = () => {
  editingCell.value = null
  editingValue.value = ''
}

// Add new translation key
const handleAddKey = async () => {
  if (!selectedProjectId.value || !newKey.value.keyName) return

  try {
    const request: BatchTranslationRequest = {
      project_id: selectedProjectId.value,
      key_name: newKey.value.keyName,
      context: newKey.value.context || undefined,
      translations: newKey.value.translations
    }

    await batchCreateTranslations(request)

    // Reset form and refresh
    newKey.value = { keyName: '', context: '', translations: {} }
    showAddKeyDialog.value = false
    await loadMatrix()
  } catch (err: any) {
    alert('添加失败: ' + (err.message || '未知错误'))
  }
}

// Delete translation key
const handleDeleteKey = async (keyName: string) => {
  if (!confirm(`确定要删除翻译键 "${keyName}" 吗？这将删除所有语言的该键翻译。`)) {
    return
  }

  if (!selectedProjectId.value) return

  try {
    // Find all translation IDs for this key
    const row = matrix.value?.rows.find((r) => r.key_name === keyName)
    if (!row) return

    const ids: number[] = []
    Object.values(row.translations).forEach((cell) => {
      if (cell.id) ids.push(cell.id)
    })

    if (ids.length > 0) {
      await api.post('/translations/batch-delete', ids)
      await loadMatrix()
    }
  } catch (err: any) {
    alert('删除失败: ' + (err.message || '未知错误'))
  }
}

// Export translations
const handleExport = async () => {
  if (!selectedProjectId.value) return

  try {
    const data = await exportTranslations(selectedProjectId.value)

    // Download as JSON file
    const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `translations-project-${selectedProjectId.value}.json`
    a.click()
    URL.revokeObjectURL(url)
  } catch (err: any) {
    alert('导出失败: ' + (err.message || '未知错误'))
  }
}

// File select for import
const handleFileSelect = (event: Event) => {
  const target = event.target as HTMLInputElement
  if (target.files && target.files.length > 0) {
    importFile.value = target.files[0]
  }
}

// Import translations
const handleImport = async () => {
  if (!selectedProjectId.value || !importFile.value) return

  try {
    const text = await importFile.value.text()
    const data: ImportTranslationsData = JSON.parse(text)

    await importTranslations(selectedProjectId.value, data)

    // Reset and refresh
    importFile.value = null
    showImportDialog.value = false
    await loadMatrix()
    alert('导入成功')
  } catch (err: any) {
    alert('导入失败: ' + (err.message || '未知错误'))
  }
}
</script>

<style scoped>
.translations-view {
  padding: 2rem;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.header h1 {
  font-size: 2rem;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.project-selector {
  padding: 0.5rem 1rem;
  border: 1px solid #d1d5db;
  border-radius: 0.375rem;
  font-size: 1rem;
  min-width: 200px;
  cursor: pointer;
}

.project-selector:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.content {
  background: white;
  border-radius: 0.5rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem;
  border-bottom: 1px solid #e5e7eb;
}

.search-box {
  flex: 1;
  max-width: 400px;
}

.search-input {
  width: 100%;
  padding: 0.5rem 1rem;
  border: 1px solid #d1d5db;
  border-radius: 0.375rem;
  font-size: 0.875rem;
}

.search-input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.toolbar-actions {
  display: flex;
  gap: 0.75rem;
}

.btn {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 0.375rem;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-primary {
  background: #3b82f6;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: #2563eb;
}

.btn-secondary {
  background: #f3f4f6;
  color: #374151;
}

.btn-secondary:hover:not(:disabled) {
  background: #e5e7eb;
}

.btn-icon {
  width: 2rem;
  height: 2rem;
  padding: 0;
  border: none;
  background: transparent;
  color: #6b7280;
  cursor: pointer;
  font-size: 1.5rem;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 0.25rem;
  transition: all 0.2s;
}

.btn-icon:hover {
  background: #f3f4f6;
  color: #374151;
}

.btn-danger:hover {
  background: #fee2e2;
  color: #dc2626;
}

.icon {
  font-size: 1rem;
}

.loading,
.error {
  padding: 3rem;
  text-align: center;
  color: #6b7280;
}

.error {
  color: #dc2626;
}

.table-wrapper {
  overflow-x: auto;
}

.translation-matrix {
  width: 100%;
  border-collapse: collapse;
}

.translation-matrix th,
.translation-matrix td {
  padding: 0.75rem 1rem;
  text-align: left;
  border-bottom: 1px solid #e5e7eb;
}

.translation-matrix th {
  background: #f9fafb;
  font-weight: 600;
  color: #374151;
  position: sticky;
  top: 0;
  z-index: 10;
}

.key-column {
  min-width: 200px;
  font-weight: 500;
}

.context-column {
  min-width: 150px;
  color: #6b7280;
  font-size: 0.875rem;
}

.language-column {
  min-width: 200px;
}

.language-code {
  margin-left: 0.5rem;
  font-size: 0.75rem;
  color: #6b7280;
  font-weight: 400;
}

.actions-column {
  width: 80px;
  text-align: center;
}

.translation-cell {
  cursor: pointer;
  transition: background 0.2s;
}

.translation-cell:hover {
  background: #f9fafb;
}

.cell-editing {
  padding: 0;
  margin: -0.75rem -1rem;
}

.cell-input {
  width: 100%;
  min-height: 60px;
  padding: 0.75rem 1rem;
  border: 2px solid #3b82f6;
  font-family: inherit;
  font-size: 0.875rem;
  resize: vertical;
}

.cell-input:focus {
  outline: none;
}

.cell-display {
  min-height: 1.5rem;
}

.cell-value {
  color: #1f2937;
}

.cell-empty {
  color: #9ca3af;
  font-style: italic;
  font-size: 0.875rem;
}

.empty-state {
  padding: 3rem !important;
  text-align: center;
  color: #6b7280;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 1rem;
  padding: 1.5rem;
  border-top: 1px solid #e5e7eb;
}

.page-info {
  color: #6b7280;
  font-size: 0.875rem;
}

.empty-project {
  padding: 4rem;
  text-align: center;
  color: #6b7280;
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
}

.modal {
  background: white;
  border-radius: 0.5rem;
  padding: 2rem;
  min-width: 400px;
  max-width: 90%;
  max-height: 85vh;
  overflow-y: auto;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1);
}

.large-modal {
  min-width: 600px;
}

.modal h2 {
  margin: 0 0 1.5rem 0;
  font-size: 1.5rem;
  color: #1f2937;
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
  color: #374151;
}

.required {
  color: #dc2626;
  margin-left: 0.25rem;
}

.form-input {
  width: 100%;
  padding: 0.5rem 0.75rem;
  border: 1px solid #d1d5db;
  border-radius: 0.375rem;
  font-size: 0.875rem;
}

.form-input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.language-inputs {
  border: 1px solid #e5e7eb;
  border-radius: 0.375rem;
  padding: 1rem;
  background: #f9fafb;
}

.language-input-row {
  display: grid;
  grid-template-columns: 180px 1fr;
  gap: 1rem;
  align-items: center;
  margin-bottom: 0.75rem;
}

.language-input-row:last-child {
  margin-bottom: 0;
}

.language-label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-weight: 500;
  color: #374151;
}

.language-name {
  font-size: 0.875rem;
}

.language-code-badge {
  display: inline-block;
  padding: 0.125rem 0.5rem;
  background: #e0e7ff;
  color: #4338ca;
  font-size: 0.75rem;
  font-weight: 600;
  border-radius: 0.25rem;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  margin-top: 1.5rem;
}
</style>
