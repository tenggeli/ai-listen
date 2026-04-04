<script setup lang="ts">
import type { ProviderSummary } from '../../domain/provider/ProviderSummary'

defineProps<{
  items: ProviderSummary[]
  activeId: string
}>()

const emit = defineEmits<{
  select: [providerId: string]
}>()
</script>

<template>
  <div class="table-wrap">
    <table>
      <thead>
        <tr>
          <th>服务方</th>
          <th>城市</th>
          <th>审核状态</th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="item in items"
          :key="item.id"
          :class="{ active: item.id === activeId }"
          @click="emit('select', item.id)"
        >
          <td>{{ item.displayName }}</td>
          <td>{{ item.cityCode }}</td>
          <td><span class="pill-soft">{{ item.reviewStatus }}</span></td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<style scoped>
.table-wrap {
  overflow: hidden;
}

table {
  width: 100%;
  border-collapse: collapse;
}

th,
td {
  text-align: left;
  padding: 12px 14px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
  font-size: 14px;
}

th {
  color: rgba(237, 247, 251, 0.5);
  font-weight: 500;
}

tr {
  cursor: pointer;
}

tr.active {
  background: rgba(115, 213, 255, 0.08);
}
</style>
