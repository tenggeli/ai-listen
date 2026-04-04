<script setup lang="ts">
import { onMounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { authService } from '../../application/auth'
import { HttpSoundAdminApi, MockSoundAdminApi, type AdminSoundItem } from '../../api/SoundAdminApi'
import AdminShell from '../../components/layout/AdminShell.vue'

const router = useRouter()
const useMock = import.meta.env.VITE_USE_MOCK === 'true'
const api = useMock
  ? new MockSoundAdminApi()
  : new HttpSoundAdminApi('/api/v1/admin', () => authService.getAccessToken())

const state = reactive<{
  loading: boolean
  saving: boolean
  actionLoading: boolean
  categoryKey: string
  status: string
  keyword: string
  pageNo: number
  pageSize: number
  items: AdminSoundItem[]
  total: number
  selectedId: string
  formMode: 'create' | 'edit'
  form: {
    id: string
    categoryKey: string
    title: string
    playCountText: string
    durationText: string
    emoji: string
    author: string
    sortOrder: number
    status: string
  }
  errorMessage: string
}>({
  loading: false,
  saving: false,
  actionLoading: false,
  categoryKey: '',
  status: '',
  keyword: '',
  pageNo: 1,
  pageSize: 20,
  items: [],
  total: 0,
  selectedId: '',
  formMode: 'create',
  form: {
    id: '',
    categoryKey: 'nature',
    title: '',
    playCountText: '',
    durationText: '',
    emoji: '',
    author: '',
    sortOrder: 0,
    status: 'inactive'
  },
  errorMessage: ''
})

onMounted(() => {
  void loadList()
})

async function loadList(): Promise<void> {
  state.loading = true
  state.errorMessage = ''
  try {
    const result = await api.list({
      categoryKey: state.categoryKey,
      status: state.status,
      keyword: state.keyword,
      pageNo: state.pageNo,
      pageSize: state.pageSize
    })
    state.items = result.items
    state.total = result.total
  } catch (error) {
    state.errorMessage = error instanceof Error ? error.message : '加载声音列表失败'
  } finally {
    state.loading = false
  }
}

function startCreate(): void {
  state.formMode = 'create'
  state.selectedId = ''
  state.form.id = ''
  state.form.categoryKey = 'nature'
  state.form.title = ''
  state.form.playCountText = ''
  state.form.durationText = ''
  state.form.emoji = ''
  state.form.author = ''
  state.form.sortOrder = 0
  state.form.status = 'inactive'
}

function startEdit(item: AdminSoundItem): void {
  state.formMode = 'edit'
  state.selectedId = item.id
  state.form.id = item.id
  state.form.categoryKey = item.categoryKey
  state.form.title = item.title
  state.form.playCountText = item.playCountText
  state.form.durationText = item.durationText
  state.form.emoji = item.emoji
  state.form.author = item.author
  state.form.sortOrder = item.sortOrder
  state.form.status = item.status
}

async function submitForm(): Promise<void> {
  state.saving = true
  state.errorMessage = ''
  try {
    if (state.formMode === 'create') {
      await api.create({
        id: state.form.id,
        categoryKey: state.form.categoryKey,
        title: state.form.title,
        playCountText: state.form.playCountText,
        durationText: state.form.durationText,
        emoji: state.form.emoji,
        author: state.form.author,
        sortOrder: state.form.sortOrder,
        status: state.form.status
      })
      startCreate()
    } else {
      await api.update(state.selectedId, {
        categoryKey: state.form.categoryKey,
        title: state.form.title,
        playCountText: state.form.playCountText,
        durationText: state.form.durationText,
        emoji: state.form.emoji,
        author: state.form.author,
        sortOrder: state.form.sortOrder
      })
    }
    await loadList()
  } catch (error) {
    state.errorMessage = error instanceof Error ? error.message : '保存失败'
  } finally {
    state.saving = false
  }
}

async function changeStatus(item: AdminSoundItem, action: 'activate' | 'deactivate'): Promise<void> {
  state.actionLoading = true
  state.errorMessage = ''
  try {
    await api.changeStatus(item.id, action)
    await loadList()
    if (state.selectedId === item.id) {
      const target = state.items.find((entry) => entry.id === item.id)
      if (target) {
        startEdit(target)
      }
    }
  } catch (error) {
    state.errorMessage = error instanceof Error ? error.message : '更新状态失败'
  } finally {
    state.actionLoading = false
  }
}

function logout(): void {
  authService.logout()
  void router.replace('/login')
}
</script>

<template>
  <AdminShell title="声音内容管理" subtitle="支持声音内容查询、编辑、创建与上/下线操作。" @logout="logout">
    <section class="admin-filters">
      <select v-model="state.categoryKey">
        <option value="">全部分类</option>
        <option value="nature">自然白噪音</option>
        <option value="sleep">睡眠引导</option>
        <option value="meditation">正念冥想</option>
        <option value="story">治愈故事</option>
        <option value="breath">呼吸练习</option>
      </select>
      <select v-model="state.status">
        <option value="">全部状态</option>
        <option value="active">active</option>
        <option value="inactive">inactive</option>
      </select>
      <input v-model.trim="state.keyword" placeholder="关键词：id/title/author" />
      <button @click="loadList">查询</button>
      <button @click="startCreate">新建</button>
    </section>

    <section class="admin-panel-grid">
      <section class="admin-list-panel admin-card">
        <p v-if="state.loading" class="admin-loading">列表加载中...</p>
        <template v-else>
          <p>共 {{ state.total }} 条</p>
          <button
            v-for="item in state.items"
            :key="item.id"
            type="button"
            class="admin-list-item"
            :class="{ active: item.id === state.selectedId }"
            @click="startEdit(item)"
          >
            <div class="admin-title-row">
              <strong>{{ item.title }}</strong>
              <span class="pill-soft">{{ item.status }}</span>
            </div>
            <p>ID: {{ item.id }}</p>
            <p>分类: {{ item.categoryKey }} · 时长: {{ item.durationText }}</p>
            <p>作者: {{ item.author || '-' }}</p>
            <div class="admin-actions">
              <button :disabled="state.actionLoading || item.status === 'active'" @click.stop="changeStatus(item, 'activate')">上线</button>
              <button :disabled="state.actionLoading || item.status === 'inactive'" @click.stop="changeStatus(item, 'deactivate')">下线</button>
            </div>
          </button>
        </template>
      </section>

      <section class="admin-card admin-detail-panel">
        <h3>{{ state.formMode === 'create' ? '新建声音内容' : `编辑声音内容：${state.selectedId}` }}</h3>
        <section class="admin-filters">
          <input v-model.trim="state.form.id" :disabled="state.formMode === 'edit'" placeholder="track_id（可选，留空自动生成）" />
          <select v-model="state.form.categoryKey">
            <option value="nature">nature</option>
            <option value="sleep">sleep</option>
            <option value="meditation">meditation</option>
            <option value="story">story</option>
            <option value="breath">breath</option>
          </select>
          <input v-model.trim="state.form.title" placeholder="标题（必填）" />
          <input v-model.trim="state.form.playCountText" placeholder="播放文案，例如：12 次播放" />
          <input v-model.trim="state.form.durationText" placeholder="时长（必填，如 10:00）" />
          <input v-model.trim="state.form.emoji" placeholder="emoji" />
          <input v-model.trim="state.form.author" placeholder="作者" />
          <input v-model.number="state.form.sortOrder" type="number" min="0" placeholder="排序" />
          <select v-if="state.formMode === 'create'" v-model="state.form.status">
            <option value="inactive">inactive</option>
            <option value="active">active</option>
          </select>
        </section>
        <button :disabled="state.saving" @click="submitForm">
          {{ state.formMode === 'create' ? '创建' : '保存更新' }}
        </button>
      </section>
    </section>

    <p v-if="state.errorMessage" class="error">{{ state.errorMessage }}</p>
  </AdminShell>
</template>
