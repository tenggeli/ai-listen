<script setup lang="ts">
import { ref } from "vue";
import { adminLogin, setAdminAccessToken } from "../../sdk";

const username = ref("admin");
const password = ref("admin123456");
const loginResult = ref("未登录");

async function handleLogin() {
  try {
    const result = await adminLogin({
      username: username.value,
      password: password.value
    });
    setAdminAccessToken(result.accessToken);
    loginResult.value = `登录成功：${result.adminUser.username}`;
    uni.showToast({ title: "登录成功", icon: "none" });
  } catch (error) {
    loginResult.value = "登录失败";
    uni.showToast({ title: "登录失败", icon: "none" });
  }
}
</script>

<template>
  <view class="page">
    <view class="title">管理员登录</view>
    <input v-model="username" class="input" placeholder="用户名" />
    <input v-model="password" class="input" password placeholder="密码" />
    <button class="btn" @click="handleLogin">登录</button>
    <view class="result">{{ loginResult }}</view>
  </view>
</template>

<style>
.page {
  min-height: 100vh;
  padding: 40rpx 32rpx;
  background: #f7fafc;
}

.title {
  margin-bottom: 32rpx;
  font-size: 40rpx;
  font-weight: 600;
  color: #15283d;
}

.input {
  height: 84rpx;
  margin-bottom: 20rpx;
  padding: 0 24rpx;
  border-radius: 12rpx;
  background: #fff;
  border: 1rpx solid #dce5ee;
}

.btn {
  border: none;
  border-radius: 12rpx;
  background: #2b7d9a;
  color: #fff;
}

.result {
  margin-top: 20rpx;
  color: #556983;
  font-size: 26rpx;
}
</style>
