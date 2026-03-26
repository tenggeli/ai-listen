<script setup lang="ts">
import { ref } from "vue";
import { designTokens, userThemeTokens } from "@listen/app-shared";
import { getHealth } from "../../sdk";
import { toUserLogin, toUserProfile } from "../../router/navigation";

const healthStatus = ref("未检查");
const tokenStyles = {
  pageBackground: userThemeTokens.pageBackground,
  cardBorder: `1rpx solid ${designTokens.color.border}`,
  primaryButton: userThemeTokens.brand,
  ghostButton: "#EBF5F9"
} as const;

async function checkHealth() {
  try {
    const data = await getHealth();
    healthStatus.value = data.status;
  } catch (error) {
    healthStatus.value = "请求失败";
    uni.showToast({ title: "健康检查失败", icon: "none" });
  }
}
</script>

<template>
  <view class="page" :style="{ background: tokenStyles.pageBackground }">
    <view class="card" :style="{ border: tokenStyles.cardBorder }">
      <view class="eyebrow">listen user app</view>
      <view class="title">用户端基础工程已迁入 App</view>
      <view class="desc">已从 user-web 迁移首页路由，后续可继续接入声音页、服务列表和订单链路。</view>
      <view class="status">服务状态：{{ healthStatus }}</view>
      <view class="actions">
        <button class="btn" :style="{ background: tokenStyles.primaryButton }" @click="checkHealth">检查后端健康</button>
        <button class="btn ghost" :style="{ background: tokenStyles.ghostButton }" @click="toUserLogin">短信登录页</button>
        <button class="btn ghost" :style="{ background: tokenStyles.ghostButton }" @click="toUserProfile">我的资料页</button>
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
  color: #2682a3;
}

.title {
  font-size: 42rpx;
  font-weight: 600;
  color: #1b2b3f;
}

.desc {
  margin-top: 20rpx;
  font-size: 28rpx;
  line-height: 1.6;
  color: #5b6e85;
}

.status {
  margin-top: 24rpx;
  font-size: 26rpx;
  color: #223a55;
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
  color: #2d85a2;
}
</style>
