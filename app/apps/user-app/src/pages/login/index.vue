<script setup lang="ts">
import { ref } from "vue";
import { loginBySMS, sendSMS, setUserAccessToken } from "../../sdk";

const mobile = ref("13800138000");
const code = ref("123456");
const loginResult = ref("未登录");

async function handleSendSMS() {
  try {
    await sendSMS({ mobile: mobile.value });
    uni.showToast({ title: "验证码已发送", icon: "none" });
  } catch (error) {
    uni.showToast({ title: "发送失败", icon: "none" });
  }
}

async function handleLogin() {
  try {
    const result = await loginBySMS({
      mobile: mobile.value,
      code: code.value
    });
    setUserAccessToken(result.accessToken);
    loginResult.value = `登录成功：${result.user.nickname}`;
    uni.showToast({ title: "登录成功", icon: "none" });
  } catch (error) {
    loginResult.value = "登录失败";
    uni.showToast({ title: "登录失败", icon: "none" });
  }
}
</script>

<template>
  <view class="page">
    <view class="title">短信登录</view>
    <input v-model="mobile" class="input" placeholder="手机号" />
    <input v-model="code" class="input" placeholder="验证码" />
    <button class="btn ghost" @click="handleSendSMS">发送验证码</button>
    <button class="btn" @click="handleLogin">登录</button>
    <view class="result">{{ loginResult }}</view>
  </view>
</template>

<style>
.page {
  min-height: 100vh;
  padding: 40rpx 32rpx;
  background: #f6f9fd;
}

.title {
  margin-bottom: 32rpx;
  font-size: 40rpx;
  font-weight: 600;
  color: #1c2f44;
}

.input {
  height: 84rpx;
  margin-bottom: 20rpx;
  padding: 0 24rpx;
  border-radius: 12rpx;
  background: #fff;
  border: 1rpx solid #dce5f0;
}

.btn {
  margin-top: 8rpx;
  border: none;
  border-radius: 12rpx;
  background: #2d85a2;
  color: #fff;
}

.btn.ghost {
  background: #eaf5f9;
  color: #2d85a2;
}

.result {
  margin-top: 20rpx;
  color: #4f627b;
  font-size: 26rpx;
}
</style>
