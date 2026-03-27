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
          <td>{{ item.reviewStatus }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<style scoped>
.table-wrap {
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  overflow: hidden;
}

table {
  width: 100%;
  border-collapse: collapse;
}

th,
td {
  text-align: left;
  padding: 10px;
  border-bottom: 1px solid #f1f5f9;
  font-size: 14px;
}

tr {
  cursor: pointer;
}

tr.active {
  background: #eff6ff;
}
</style>
