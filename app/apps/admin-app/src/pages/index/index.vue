<script setup lang="ts">
import { ref } from "vue";
import { adminThemeTokens, designTokens } from "@listen/app-shared";
import { adminMe, dashboardOverview } from "../../sdk";
import { toAdminLogin } from "../../router/navigation";

const overviewText = ref("未加载");
const meText = ref("未加载");
const tokenStyles = {
  pageBackground: adminThemeTokens.pageBackground,
  cardBorder: `1rpx solid ${designTokens.color.border}`,
  primaryButton: adminThemeTokens.brand,
  ghostButton: "#EAF4F8"
} as const;

async function loadOverview() {
  try {
    const data = await dashboardOverview();
    overviewText.value = `${data.module}:${data.action}`;
  } catch (error) {
    overviewText.value = "加载失败";
    uni.showToast({ title: "看板加载失败", icon: "none" });
  }
}

async function loadMe() {
  try {
    const data = await adminMe();
    meText.value = `管理员：${data.adminUser.username}`;
  } catch (error) {
    meText.value = "加载失败，请先登录";
    uni.showToast({ title: "请先登录", icon: "none" });
  }
}
</script>

<template>
  <view class="page" :style="{ background: tokenStyles.pageBackground }">
    <view class="card" :style="{ border: tokenStyles.cardBorder }">
      <view class="eyebrow">listen admin app</view>
      <view class="title">管理端基础工程已迁入 App</view>
      <view class="desc">已从 admin-web 迁移控制台路由，后续可继续接入审核、订单、财务和配置中心。</view>
      <view class="status">看板接口：{{ overviewText }}</view>
      <view class="status">当前管理员：{{ meText }}</view>
      <view class="actions">
        <button class="btn" :style="{ background: tokenStyles.primaryButton }" @click="loadOverview">加载看板</button>
        <button class="btn ghost" :style="{ background: tokenStyles.ghostButton }" @click="loadMe">加载当前管理员</button>
        <button class="btn ghost" :style="{ background: tokenStyles.ghostButton }" @click="toAdminLogin">前往登录页</button>
      </view>
    </view>
  </view>
</template>

<style>
.page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 32rpx;
}

.card {
  width: 100%;
  max-width: 680rpx;
  border-radius: 24rpx;
  padding: 36rpx;
  background: #ffffff;
}

.eyebrow {
  margin-bottom: 14rpx;
  font-size: 22rpx;
  letter-spacing: 4rpx;
  text-transform: uppercase;
  color: #2b7d9a;
}

.title {
  font-size: 42rpx;
  font-weight: 600;
  color: #14263a;
}

.desc {
  margin-top: 20rpx;
  font-size: 28rpx;
  line-height: 1.6;
  color: #5b6e85;
}

.status {
  margin-top: 20rpx;
  font-size: 26rpx;
  color: #21384f;
}

.actions {
  margin-top: 24rpx;
  display: flex;
  flex-direction: column;
  gap: 16rpx;
}

.btn {
  border: none;
  border-radius: 14rpx;
  color: #fff;
}

.btn.ghost {
  color: #2b7d9a;
}
</style>
