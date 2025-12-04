<template>
  <div class="page">
    <header class="page-head">
      <div>
        <h2>分类与标签</h2>
        <p class="subtitle">维护平台展示结构，保障前台筛选体验。</p>
      </div>
      <div class="actions">
        <button
          type="button"
          class="btn primary"
          :disabled="!canManage"
          @click="openCategoryModal('create')"
        >
          新增分类
        </button>
        <button
          type="button"
          class="btn primary"
          :disabled="!canManage"
          @click="openTagModal('create')"
        >
          新增标签
        </button>
      </div>
    </header>

    <section class="section">
      <header class="section-head">
        <h3>商品分类</h3>
        <span class="section-desc">共 {{ categoryList.length }} 个</span>
      </header>
      <ul class="list">
        <li v-for="item in categoryList" :key="item.id" class="list-item">
          <div>
            <p class="item-title">{{ item.name }}</p>
            <p class="item-desc">{{ item.description }}</p>
          </div>
          <div class="item-actions">
            <span class="item-count">关联商品：{{ item.productCount }}</span>
            <button
              type="button"
              class="link"
              :disabled="!canManage"
              @click="openCategoryModal('edit', item.id)"
            >
              编辑
            </button>
            <button
              type="button"
              class="link danger"
              :disabled="!canManage"
              @click="openDeleteModal('category', item.id)"
            >
              删除
            </button>
          </div>
        </li>
      </ul>
    </section>

    <section class="section">
      <header class="section-head">
        <h3>热门标签</h3>
        <span class="section-desc">共 {{ tagList.length }} 个</span>
      </header>
      <div class="tag-grid">
        <article v-for="tag in tagList" :key="tag.id" class="tag-card">
          <h4>{{ tag.name }}</h4>
          <p>{{ tag.description }}</p>
          <footer>
            <span>使用次数：{{ tag.usage }}</span>
            <div class="tag-actions">
              <button
                type="button"
                class="link"
                :disabled="!canManage"
                @click="openTagModal('edit', tag.id)"
              >
                编辑
              </button>
              <button
                type="button"
                class="link danger"
                :disabled="!canManage"
                @click="openDeleteModal('tag', tag.id)"
              >
                删除
              </button>
            </div>
          </footer>
        </article>
      </div>
    </section>
  </div>

  <BaseModal
    :open="categoryModal.open"
    :title="categoryModal.mode === 'create' ? '新增分类' : '编辑分类'"
    width="520px"
    @close="closeCategoryModal"
  >
    <form class="modal-form" @submit.prevent="submitCategoryModal">
      <div class="modal-field">
        <label class="modal-label" for="category-name-input">分类名称</label>
        <input
          id="category-name-input"
          v-model="categoryModal.name"
          class="modal-input"
          placeholder="请输入分类名称"
        />
      </div>
      <div class="modal-field full">
        <label class="modal-label" for="category-desc-input">分类描述</label>
        <textarea
          id="category-desc-input"
          v-model="categoryModal.description"
          class="modal-textarea"
          rows="3"
          placeholder="用于描述该分类的主要内容"
        ></textarea>
      </div>
      <p v-if="categoryModal.error" class="modal-error">{{ categoryModal.error }}</p>
      <div class="modal-actions">
        <button type="button" class="btn ghost" @click="closeCategoryModal">取消</button>
        <button type="submit" class="btn primary" :disabled="!canManage">
          {{ categoryModal.mode === 'create' ? '创建' : '保存' }}
        </button>
      </div>
    </form>
  </BaseModal>

  <BaseModal
    :open="tagModal.open"
    :title="tagModal.mode === 'create' ? '新增标签' : '编辑标签'"
    width="520px"
    @close="closeTagModal"
  >
    <form class="modal-form" @submit.prevent="submitTagModal">
      <div class="modal-field">
        <label class="modal-label" for="tag-name-input">标签名称</label>
        <input
          id="tag-name-input"
          v-model="tagModal.name"
          class="modal-input"
          placeholder="请输入标签名称"
        />
      </div>
      <div class="modal-field full">
        <label class="modal-label" for="tag-desc-input">标签描述</label>
        <textarea
          id="tag-desc-input"
          v-model="tagModal.description"
          class="modal-textarea"
          rows="3"
          placeholder="用于商品筛选或展示的说明"
        ></textarea>
      </div>
      <p v-if="tagModal.error" class="modal-error">{{ tagModal.error }}</p>
      <div class="modal-actions">
        <button type="button" class="btn ghost" @click="closeTagModal">取消</button>
        <button type="submit" class="btn primary" :disabled="!canManage">
          {{ tagModal.mode === 'create' ? '创建' : '保存' }}
        </button>
      </div>
    </form>
  </BaseModal>

  <BaseModal
    :open="Boolean(deleteModal)"
    title="确认删除"
    width="420px"
    @close="closeDeleteModal"
  >
    <template v-if="deleteModal">
      <p class="modal-message">
        {{
          deleteModal.type === 'category'
            ? `删除分类 ${deleteModal.name} 后将无法恢复，确定继续吗？`
            : `删除标签 ${deleteModal.name} 后将无法恢复，确定继续吗？`
        }}
      </p>
      <div class="modal-actions">
        <button type="button" class="btn ghost" @click="closeDeleteModal">取消</button>
        <button type="button" class="btn primary" :disabled="!canManage" @click="submitDeleteModal">
          确认删除
        </button>
      </div>
    </template>
  </BaseModal>

  <BaseModal
    :open="Boolean(feedbackModal)"
    title="操作提示"
    width="420px"
    @close="closeFeedback"
  >
    <p class="modal-message">{{ feedbackModal?.message }}</p>
    <div class="modal-actions">
      <button type="button" class="btn primary" @click="closeFeedback">知道了</button>
    </div>
  </BaseModal>
</template>

<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import { storeToRefs } from 'pinia'

import BaseModal from '../components/BaseModal.vue'
import { useAdminDataStore } from '../stores/admin'
import { useAuthStore } from '../stores/auth'
import type { AdminCategory, AdminTag } from '../stores/admin'

defineOptions({ name: 'AdminCategoryManageView' })

const adminStore = useAdminDataStore()
const authStore = useAuthStore()
const { categories, tags, productCountsByCategory } = storeToRefs(adminStore)
const { canManage } = storeToRefs(authStore)

interface AdminCategoryWithCount extends AdminCategory {
  productCount: number
}

const categoryList = computed<AdminCategoryWithCount[]>(() =>
  categories.value.map((item: AdminCategory) => ({
    ...item,
    productCount: productCountsByCategory.value.get(item.id) ?? 0,
  })),
)

const tagList = computed<AdminTag[]>(() => tags.value)

const feedbackModal = ref<{ message: string } | null>(null)

const showPermissionWarning = () => {
  feedbackModal.value = { message: '当前账号仅支持查看数据，如需操作请联系管理员。' }
}

type DialogMode = 'create' | 'edit'

const categoryModal = reactive({
  open: false,
  mode: 'create' as DialogMode,
  id: null as number | null,
  name: '',
  description: '',
  error: '',
})

const tagModal = reactive({
  open: false,
  mode: 'create' as DialogMode,
  id: null as number | null,
  name: '',
  description: '',
  error: '',
})

const deleteModal = ref<{ type: 'category' | 'tag'; id: number; name: string } | null>(null)

const openCategoryModal = (mode: DialogMode, id?: number) => {
  if (!canManage.value) {
    showPermissionWarning()
    return
  }
  categoryModal.open = false
  categoryModal.mode = mode
  categoryModal.error = ''
  if (mode === 'edit') {
    if (typeof id !== 'number') return
    const category = categories.value.find((item) => item.id === id)
    if (!category) {
      feedbackModal.value = { message: '未找到该分类。' }
      return
    }
    categoryModal.id = category.id
    categoryModal.name = category.name
    categoryModal.description = category.description
  } else {
    categoryModal.id = null
    categoryModal.name = ''
    categoryModal.description = '用于描述该分类的主要内容'
  }
  categoryModal.open = true
}

const closeCategoryModal = () => {
  categoryModal.open = false
}

const submitCategoryModal = () => {
  if (!canManage.value) {
    showPermissionWarning()
    categoryModal.open = false
    return
  }
  const name = categoryModal.name.trim()
  if (!name) {
    categoryModal.error = '分类名称不能为空。'
    return
  }
  const description = categoryModal.description.trim()

  if (categoryModal.mode === 'create') {
    adminStore.addCategory({ name, description })
    categoryModal.open = false
    feedbackModal.value = { message: '分类已创建。' }
    return
  }

  if (categoryModal.id === null) {
    categoryModal.error = '未找到该分类。'
    return
  }

  const success = adminStore.updateCategory(categoryModal.id, { name, description })
  if (!success) {
    categoryModal.error = '更新分类时出现问题。'
    return
  }
  categoryModal.open = false
  feedbackModal.value = { message: '分类信息已更新。' }
}

const openTagModal = (mode: DialogMode, id?: number) => {
  if (!canManage.value) {
    showPermissionWarning()
    return
  }
  tagModal.open = false
  tagModal.mode = mode
  tagModal.error = ''
  if (mode === 'edit') {
    if (typeof id !== 'number') return
    const tag = tags.value.find((item) => item.id === id)
    if (!tag) {
      feedbackModal.value = { message: '未找到该标签。' }
      return
    }
    tagModal.id = tag.id
    tagModal.name = tag.name
    tagModal.description = tag.description
  } else {
    tagModal.id = null
    tagModal.name = ''
    tagModal.description = '用于商品筛选或展示'
  }
  tagModal.open = true
}

const closeTagModal = () => {
  tagModal.open = false
}

const submitTagModal = () => {
  if (!canManage.value) {
    showPermissionWarning()
    tagModal.open = false
    return
  }
  const name = tagModal.name.trim()
  if (!name) {
    tagModal.error = '标签名称不能为空。'
    return
  }
  const description = tagModal.description.trim()

  if (tagModal.mode === 'create') {
    adminStore.addTag({ name, description })
    tagModal.open = false
    feedbackModal.value = { message: '标签已创建。' }
    return
  }

  if (tagModal.id === null) {
    tagModal.error = '未找到该标签。'
    return
  }

  const success = adminStore.updateTag(tagModal.id, { name, description })
  if (!success) {
    tagModal.error = '更新标签时出现问题。'
    return
  }
  tagModal.open = false
  feedbackModal.value = { message: '标签已更新。' }
}

const openDeleteModal = (type: 'category' | 'tag', id: number) => {
  if (!canManage.value) {
    showPermissionWarning()
    return
  }
  if (type === 'category') {
    const category = categories.value.find((item) => item.id === id)
    if (!category) {
      feedbackModal.value = { message: '未找到该分类。' }
      return
    }
    const count = productCountsByCategory.value.get(id) ?? 0
    if (count > 0) {
      feedbackModal.value = { message: '该分类仍有关联商品，无法删除。' }
      return
    }
    deleteModal.value = { type, id, name: category.name }
    return
  }

  const tag = tags.value.find((item) => item.id === id)
  if (!tag) {
    feedbackModal.value = { message: '未找到该标签。' }
    return
  }
  deleteModal.value = { type, id, name: tag.name }
}

const closeDeleteModal = () => {
  deleteModal.value = null
}

const submitDeleteModal = () => {
  const modal = deleteModal.value
  if (!modal) return

  if (!canManage.value) {
    showPermissionWarning()
    deleteModal.value = null
    return
  }

  if (modal.type === 'category') {
    const success = adminStore.removeCategory(modal.id)
    deleteModal.value = null
    feedbackModal.value = { message: success ? '分类已删除。' : '该分类仍有关联商品，无法删除。' }
    return
  }

  const success = adminStore.removeTag(modal.id)
  deleteModal.value = null
  feedbackModal.value = { message: success ? '标签已删除。' : '删除标签失败。' }
}

const closeFeedback = () => {
  feedbackModal.value = null
}
</script>

<style scoped>
.page {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.page-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.page-head h2 {
  margin: 0;
}

.subtitle {
  margin: 6px 0 0;
  color: #6b7280;
}

.actions {
  display: flex;
  gap: 12px;
}

.btn {
  border: 1px solid #2563eb;
  background: transparent;
  color: #2563eb;
  padding: 8px 16px;
  border-radius: 6px;
  cursor: pointer;
}

.btn.primary {
  background: #2563eb;
  color: #fff;
  border-color: #2563eb;
}

.btn:disabled {
  cursor: not-allowed;
  opacity: 0.55;
  border-color: #d1d5db;
  color: #9ca3af;
  background: #f3f4f6;
}

.btn.primary:disabled {
  background: #93c5fd;
  color: #fff;
}

.section {
  background: #fff;
  border-radius: 12px;
  padding: 20px;
  border: 1px solid #e5e7eb;
  box-shadow: 0 4px 16px rgba(15, 23, 42, 0.08);
}

.section-head {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  margin-bottom: 12px;
}

.section-head h3 {
  margin: 0;
}

.section-desc {
  color: #6b7280;
  font-size: 13px;
}

.list {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.list-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #f9fafb;
  padding: 16px;
  border-radius: 10px;
}

.item-title {
  margin: 0 0 6px;
  font-weight: 600;
}

.item-desc {
  margin: 0;
  color: #6b7280;
  font-size: 13px;
}

.item-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.item-count {
  color: #6b7280;
  font-size: 13px;
}

.tag-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 16px;
}

.tag-card {
  background: #f9fafb;
  padding: 18px;
  border-radius: 12px;
  border: 1px solid #e5e7eb;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.tag-card h4 {
  margin: 0;
}

.tag-card p {
  margin: 0;
  color: #6b7280;
  font-size: 13px;
}

.tag-card footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 13px;
  color: #6b7280;
}

.tag-actions {
  display: flex;
  gap: 10px;
}

.link {
  border: none;
  background: transparent;
  color: #2563eb;
  cursor: pointer;
  padding: 0;
}

.link.danger {
  color: #dc2626;
}

.link:disabled {
  cursor: not-allowed;
  color: #9ca3af;
}

.link.danger:disabled {
  color: #fca5a5;
}

.btn.ghost {
  background: transparent;
  border: 1px solid transparent;
  color: #4b5563;
}

.modal-form {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 16px;
}

.modal-field {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.modal-field.full {
  grid-column: 1 / -1;
}

.modal-label {
  font-size: 14px;
  font-weight: 500;
  color: #374151;
}

.modal-input,
.modal-textarea {
  width: 100%;
  padding: 10px 12px;
  border-radius: 8px;
  border: 1px solid #d1d5db;
  font-size: 14px;
  transition: border-color 0.2s ease, box-shadow 0.2s ease;
}

.modal-input:focus,
.modal-textarea:focus {
  border-color: #2563eb;
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.12);
  outline: none;
}

.modal-textarea {
  resize: vertical;
  min-height: 96px;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  grid-column: 1 / -1;
  margin-top: 4px;
}

.modal-error {
  margin: -4px 0 0;
  color: #dc2626;
  font-size: 13px;
  grid-column: 1 / -1;
}

.modal-message {
  margin: 0 0 12px;
  font-size: 15px;
  color: #374151;
}
</style>
