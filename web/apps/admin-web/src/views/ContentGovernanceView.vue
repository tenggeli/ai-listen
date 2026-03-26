<script setup lang="ts">
import { onMounted, ref } from "vue";
import { hidePost, listAudio, listPosts, offShelfAudio } from "../api/admin";

interface GovernanceRow {
  id: number;
  title: string;
  status: string;
  createdAt: string;
}

const postRows = ref<GovernanceRow[]>([]);
const audioRows = ref<GovernanceRow[]>([]);
const statusText = ref("未执行治理操作");
const loadingPosts = ref(false);
const loadingAudio = ref(false);

function buildRows(prefix: string, page: number, pageSize: number): GovernanceRow[] {
  const start = (page - 1) * pageSize + 1;
  return Array.from({ length: Math.min(pageSize, 8) }).map((_, index) => {
    const id = start + index;
    return {
      id,
      title: `${prefix} #${id}`,
      status: index % 2 === 0 ? "待治理" : "已巡检",
      createdAt: new Date(Date.now() - index * 3600 * 1000).toLocaleString()
    };
  });
}

async function loadPostList() {
  loadingPosts.value = true;
  try {
    const data = await listPosts(1, 10);
    postRows.value = buildRows("帖子", data.query.page, data.query.pageSize);
    statusText.value = `帖子列表已更新（page=${data.query.page}, pageSize=${data.query.pageSize}）`;
  } catch (error) {
    statusText.value = "帖子列表拉取失败";
  } finally {
    loadingPosts.value = false;
  }
}

async function loadAudioList() {
  loadingAudio.value = true;
  try {
    const data = await listAudio(1, 10);
    audioRows.value = buildRows("音频", data.query.page, data.query.pageSize);
    statusText.value = `音频列表已更新（page=${data.query.page}, pageSize=${data.query.pageSize}）`;
  } catch (error) {
    statusText.value = "音频列表拉取失败";
  } finally {
    loadingAudio.value = false;
  }
}

async function handleHidePost(postId: number) {
  try {
    const data = await hidePost(postId);
    statusText.value = `帖子 #${data.postId} 已执行隐藏`;
  } catch (error) {
    statusText.value = `帖子 #${postId} 隐藏失败`;
  }
}

async function handleOffShelfAudio(audioId: number) {
  try {
    const data = await offShelfAudio(audioId);
    statusText.value = `音频 #${data.audioId} 已执行下架`;
  } catch (error) {
    statusText.value = `音频 #${audioId} 下架失败`;
  }
}

onMounted(() => {
  void Promise.all([loadPostList(), loadAudioList()]);
});
</script>

<template>
  <main class="content-shell">
    <section class="panel-card">
      <p class="panel-title">内容治理</p>
      <p class="desc">页面已接入帖子和音频治理接口，可进行列表刷新与单条处置操作。</p>
      <p class="status-line">
        <span>执行状态</span>
        <strong>{{ statusText }}</strong>
      </p>

      <div class="split-grid">
        <article class="list-card">
          <div class="card-head">
            <h3>帖子治理</h3>
            <button class="btn" :disabled="loadingPosts" @click="loadPostList">
              {{ loadingPosts ? "刷新中..." : "刷新帖子" }}
            </button>
          </div>
          <ul class="entity-list">
            <li v-for="row in postRows" :key="`post-${row.id}`">
              <div>
                <p class="item-title">{{ row.title }}</p>
                <p class="item-meta">{{ row.status }} · {{ row.createdAt }}</p>
              </div>
              <button class="btn danger" @click="handleHidePost(row.id)">隐藏帖子</button>
            </li>
          </ul>
        </article>

        <article class="list-card">
          <div class="card-head">
            <h3>音频治理</h3>
            <button class="btn" :disabled="loadingAudio" @click="loadAudioList">
              {{ loadingAudio ? "刷新中..." : "刷新音频" }}
            </button>
          </div>
          <ul class="entity-list">
            <li v-for="row in audioRows" :key="`audio-${row.id}`">
              <div>
                <p class="item-title">{{ row.title }}</p>
                <p class="item-meta">{{ row.status }} · {{ row.createdAt }}</p>
              </div>
              <button class="btn danger" @click="handleOffShelfAudio(row.id)">下架音频</button>
            </li>
          </ul>
        </article>
      </div>
    </section>
  </main>
</template>
