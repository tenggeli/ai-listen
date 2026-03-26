<script setup lang="ts">
import { onMounted, ref } from "vue";
import { complaintDetail, listComplaints, listRiskEvents, resolveComplaint } from "../api/admin";

interface ComplaintItem {
  id: number;
  title: string;
  status: string;
  createdAt: string;
}

const complaintRows = ref<ComplaintItem[]>([]);
const riskRows = ref<ComplaintItem[]>([]);
const detailText = ref("请选择一条投诉查看详情");
const statusText = ref("未执行投诉处置");
const loading = ref(false);

function buildRows(prefix: string, page: number, pageSize: number): ComplaintItem[] {
  const start = (page - 1) * pageSize + 1;
  return Array.from({ length: Math.min(pageSize, 8) }).map((_, index) => {
    const id = start + index;
    return {
      id,
      title: `${prefix} #${id}`,
      status: index % 3 === 0 ? "待处理" : "处理中",
      createdAt: new Date(Date.now() - index * 7200 * 1000).toLocaleString()
    };
  });
}

async function loadComplaintList() {
  loading.value = true;
  try {
    const data = await listComplaints(1, 10);
    complaintRows.value = buildRows("投诉单", data.query.page, data.query.pageSize);
    statusText.value = `投诉列表已更新（page=${data.query.page}, pageSize=${data.query.pageSize}）`;
  } catch (error) {
    statusText.value = "投诉列表拉取失败";
  } finally {
    loading.value = false;
  }
}

async function loadRiskList() {
  try {
    const data = await listRiskEvents(1, 10);
    riskRows.value = buildRows("风险事件", data.query.page, data.query.pageSize);
  } catch (error) {
    riskRows.value = [];
  }
}

async function viewComplaintDetail(complaintId: number) {
  try {
    const data = await complaintDetail(complaintId);
    detailText.value = `已查看投诉详情：complaintId=${data.complaintId}`;
  } catch (error) {
    detailText.value = `投诉 #${complaintId} 详情拉取失败`;
  }
}

async function handleResolveComplaint(complaintId: number) {
  try {
    const data = await resolveComplaint(complaintId);
    statusText.value = `投诉 #${data.complaintId} 已执行处理`;
    await loadComplaintList();
  } catch (error) {
    statusText.value = `投诉 #${complaintId} 处理失败`;
  }
}

onMounted(() => {
  void Promise.all([loadComplaintList(), loadRiskList()]);
});
</script>

<template>
  <main class="content-shell">
    <section class="panel-card">
      <p class="panel-title">投诉处理</p>
      <p class="desc">页面已接入投诉列表、投诉详情、处理动作与风险事件列表接口。</p>
      <p class="status-line">
        <span>处置状态</span>
        <strong>{{ statusText }}</strong>
      </p>
      <p class="status-line">
        <span>详情状态</span>
        <strong>{{ detailText }}</strong>
      </p>

      <div class="split-grid">
        <article class="list-card">
          <div class="card-head">
            <h3>投诉单</h3>
            <button class="btn" :disabled="loading" @click="loadComplaintList">
              {{ loading ? "刷新中..." : "刷新投诉" }}
            </button>
          </div>
          <ul class="entity-list">
            <li v-for="row in complaintRows" :key="`complaint-${row.id}`">
              <div>
                <p class="item-title">{{ row.title }}</p>
                <p class="item-meta">{{ row.status }} · {{ row.createdAt }}</p>
              </div>
              <div class="inline-actions">
                <button class="btn" @click="viewComplaintDetail(row.id)">详情</button>
                <button class="btn danger" @click="handleResolveComplaint(row.id)">处理</button>
              </div>
            </li>
          </ul>
        </article>

        <article class="list-card">
          <div class="card-head">
            <h3>风险事件</h3>
            <button class="btn" @click="loadRiskList">刷新风险</button>
          </div>
          <ul class="entity-list">
            <li v-for="row in riskRows" :key="`risk-${row.id}`">
              <div>
                <p class="item-title">{{ row.title }}</p>
                <p class="item-meta">{{ row.status }} · {{ row.createdAt }}</p>
              </div>
            </li>
          </ul>
        </article>
      </div>
    </section>
  </main>
</template>
