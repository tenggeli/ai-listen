<script setup lang="ts">
import { onMounted, ref } from "vue";
import { listComplaints, listPosts, listRiskEvents } from "../api/admin";

const postSummary = ref("未拉取");
const complaintSummary = ref("未拉取");
const riskSummary = ref("未拉取");

async function loadOverview() {
  try {
    const [posts, complaints, risks] = await Promise.all([
      listPosts(1, 10),
      listComplaints(1, 10),
      listRiskEvents(1, 10)
    ]);
    postSummary.value = `分页参数 page=${posts.query.page}, pageSize=${posts.query.pageSize}`;
    complaintSummary.value = `分页参数 page=${complaints.query.page}, pageSize=${complaints.query.pageSize}`;
    riskSummary.value = `分页参数 page=${risks.query.page}, pageSize=${risks.query.pageSize}`;
  } catch (error) {
    postSummary.value = "请求失败";
    complaintSummary.value = "请求失败";
    riskSummary.value = "请求失败";
  }
}

onMounted(() => {
  void loadOverview();
});
</script>

<template>
  <main class="content-shell">
    <section class="panel-card">
      <p class="panel-title">治理看板</p>
      <p class="desc">已补齐治理入口路由，可直接进入内容治理与投诉处理页面执行处置操作。</p>
      <div class="kpi-grid">
        <article class="kpi-card">
          <h3>帖子治理</h3>
          <p>{{ postSummary }}</p>
        </article>
        <article class="kpi-card">
          <h3>投诉单列表</h3>
          <p>{{ complaintSummary }}</p>
        </article>
        <article class="kpi-card">
          <h3>风险事件</h3>
          <p>{{ riskSummary }}</p>
        </article>
      </div>
      <div class="actions">
        <button class="btn" @click="loadOverview">刷新看板数据</button>
      </div>
    </section>
  </main>
</template>
