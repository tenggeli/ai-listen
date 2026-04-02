<script setup lang="ts">
import { useRouter } from 'vue-router'
import { authService } from '../../application/auth'

const router = useRouter()

function logout(): void {
  authService.logout()
  void router.replace('/login')
}
</script>

<template>
  <main class="page">
    <div class="top-nav">
      <span>欢迎进入后台</span>
      <button @click="logout">退出登录</button>
    </div>
    <h1>管理后台骨架</h1>
    <p>当前已接入服务方审核与服务项目管理模块。</p>
    <p><RouterLink to="/providers/review">进入服务方审核</RouterLink></p>
    <p><RouterLink to="/services/manage">进入服务项目管理</RouterLink></p>
  </main>
</template>

<style scoped>
.page {
  padding: 24px;
}

.top-nav {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
</style>
