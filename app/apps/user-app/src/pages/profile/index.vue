<script setup lang="ts">
import { ref } from "vue";
import { getMe } from "../../sdk";

const profileText = ref("点击按钮加载");

async function loadProfile() {
  try {
    const data = await getMe();
    profileText.value = `昵称：${data.user.nickname}，手机号：${data.user.mobile}`;
  } catch (error) {
    profileText.value = "加载失败，请先登录";
    uni.showToast({ title: "加载失败", icon: "none" });
  }
}
</script>

<template>
  <view class="page">
    <view class="title">我的资料</view>
    <button class="btn" @click="loadProfile">加载我的信息</button>
    <view class="result">{{ profileText }}</view>
  </view>
</template>

<style>
.page {
  min-height: 100vh;
  padding: 40rpx 32rpx;
  background: #f6f9fd;
}

.title {
  margin-bottom: 24rpx;
  font-size: 40rpx;
  font-weight: 600;
  color: #1c2f44;
}

.btn {
  border: none;
  border-radius: 12rpx;
  background: #2d85a2;
  color: #fff;
}

.result {
  margin-top: 24rpx;
  color: #4f627b;
  font-size: 28rpx;
  line-height: 1.6;
}
</style>
