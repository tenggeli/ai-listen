<script setup lang="ts">
import { onMounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { authService } from '../../application/auth'
import { HttpSoundAdminApi, MockSoundAdminApi, type AdminSoundItem } from '../../api/SoundAdminApi'

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
  <main class="page">
    <nav class="top-nav">
      <RouterLink to="/dashboard">返回仪表盘</RouterLink>
      <button @click="logout">退出登录</button>
    </nav>

    <header>
      <h1>声音内容管理</h1>
      <p>首版能力：列表、创建、编辑、上/下线。</p>
    </header>

    <section class="filters">
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

    <section class="layout">
      <section class="list-panel">
        <p v-if="state.loading">列表加载中...</p>
        <template v-else>
          <p class="total">共 {{ state.total }} 条</p>
          <button
            v-for="item in state.items"
            :key="item.id"
            type="button"
            class="list-item"
            :class="{ active: item.id === state.selectedId }"
            @click="startEdit(item)"
          >
            <div class="title-row">
              <strong>{{ item.title }}</strong>
              <span class="status">{{ item.status }}</span>
            </div>
            <p>ID: {{ item.id }}</p>
            <p>分类: {{ item.categoryKey }} · 时长: {{ item.durationText }}</p>
            <p>作者: {{ item.author || '-' }}</p>
            <div class="actions">
              <button :disabled="state.actionLoading || item.status === 'active'" @click.stop="changeStatus(item, 'activate')">上线</button>
              <button :disabled="state.actionLoading || item.status === 'inactive'" @click.stop="changeStatus(item, 'deactivate')">下线</button>
            </div>
          </button>
        </template>
      </section>

      <section class="form-panel">
        <h3>{{ state.formMode === 'create' ? '新建声音内容' : `编辑声音内容：${state.selectedId}` }}</h3>
        <div class="form-grid">
          <label>
            <span>track_id（可选）</span>
            <input v-model.trim="state.form.id" :disabled="state.formMode === 'edit'" placeholder="留空自动生成" />
          </label>
          <label>
            <span>分类</span>
            <select v-model="state.form.categoryKey">
              <option value="nature">nature</option>
              <option value="sleep">sleep</option>
              <option value="meditation">meditation</option>
              <option value="story">story</option>
              <option value="breath">breath</option>
            </select>
          </label>
          <label>
            <span>标题</span>
            <input v-model.trim="state.form.title" placeholder="必填" />
          </label>
          <label>
            <span>播放文案</span>
            <input v-model.trim="state.form.playCountText" placeholder="例如：12 次播放" />
          </label>
          <label>
            <span>时长</span>
            <input v-model.trim="state.form.durationText" placeholder="必填，如 10:00" />
          </label>
          <label>
            <span>emoji</span>
            <input v-model.trim="state.form.emoji" placeholder="例如：🌙" />
          </label>
          <label>
            <span>作者</span>
            <input v-model.trim="state.form.author" placeholder="例如：listen 声音库" />
          </label>
          <label>
            <span>排序</span>
            <input v-model.number="state.form.sortOrder" type="number" min="0" />
          </label>
          <label v-if="state.formMode === 'create'">
            <span>初始状态</span>
            <select v-model="state.form.status">
              <option value="inactive">inactive</option>
              <option value="active">active</option>
            </select>
          </label>
        </div>
        <button class="submit-btn" :disabled="state.saving" @click="submitForm">
          {{ state.formMode === 'create' ? '创建' : '保存更新' }}
        </button>
      </section>
    </section>

    <p v-if="state.errorMessage" class="error">{{ state.errorMessage }}</p>
  </main>
</template>

<style scoped>
.page {
  padding: 20px;
}

.top-nav {
  margin-bottom: 12px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.filters {
  display: grid;
  grid-template-columns: 1fr 1fr 2fr auto auto;
  gap: 8px;
  margin: 12px 0;
}

.filters select,
.filters input,
.filters button,
.form-grid input,
.form-grid select,
.submit-btn {
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  padding: 8px 10px;
  background: #fff;
}

.layout {
  display: grid;
  grid-template-columns: 1.1fr 1fr;
  gap: 16px;
}

.list-panel {
  display: grid;
  gap: 8px;
}

.total {
  margin: 0 0 4px;
  color: #475569;
}

.list-item {
  text-align: left;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  background: #fff;
  padding: 10px;
}

.list-item.active {
  border-color: #0f172a;
}

.title-row {
  display: flex;
  justify-content: space-between;
  gap: 8px;
}

.status {
  color: #334155;
  font-size: 12px;
}

.actions {
  display: flex;
  gap: 8px;
}

.actions button {
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  padding: 6px 10px;
  background: #fff;
}

.form-panel {
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  padding: 12px;
  background: #fff;
}

.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
}

.form-grid label {
  display: grid;
  gap: 6px;
  font-size: 13px;
  color: #334155;
}

.submit-btn {
  margin-top: 12px;
}

.error {
  margin-top: 10px;
  color: #b91c1c;
}

@media (max-width: 960px) {
  .filters {
    grid-template-columns: 1fr 1fr;
  }

  .layout {
    grid-template-columns: 1fr;
  }

  .form-grid {
    grid-template-columns: 1fr;
  }
}
</style>
