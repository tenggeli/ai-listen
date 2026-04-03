import type { ProviderOrder, ProviderOrderStatus } from '../../domain/order/ProviderOrder'

export type OrderStatusFilter = 'all' | ProviderOrderStatus
export type ProviderOrderAction = 'accept' | 'depart' | 'arrive' | 'start' | 'complete'

export function getOrderStatusLabel(status: ProviderOrderStatus): string {
  switch (status) {
    case 'paid':
      return '待接单'
    case 'accepted':
      return '已接单'
    case 'on_the_way':
      return '出发中'
    case 'arrived':
      return '已到达'
    case 'in_service':
      return '服务中'
    case 'completed':
      return '已完单'
    case 'created':
    default:
      return '待支付'
  }
}

export function getOrderStatusTagType(status: ProviderOrderStatus): 'warn' | 'info' | 'success' | 'default' {
  switch (status) {
    case 'paid':
      return 'warn'
    case 'accepted':
    case 'on_the_way':
    case 'arrived':
      return 'info'
    case 'in_service':
    case 'completed':
      return 'success'
    default:
      return 'default'
  }
}

export function formatOrderTime(time: string | null): string {
  if (!time) {
    return '暂无'
  }
  const value = new Date(time)
  if (Number.isNaN(value.getTime())) {
    return '暂无'
  }
  return value.toLocaleString('zh-CN')
}

export function filterOrdersByStatus(items: ProviderOrder[], status: OrderStatusFilter): ProviderOrder[] {
  if (status === 'all') {
    return items
  }
  return items.filter((item) => item.status === status)
}

export function getTotalPages(total: number, pageSize: number): number {
  if (total <= 0 || pageSize <= 0) {
    return 1
  }
  return Math.max(1, Math.ceil(total / pageSize))
}

export function getNextOrderAction(status: ProviderOrderStatus): ProviderOrderAction | null {
  switch (status) {
    case 'paid':
      return 'accept'
    case 'accepted':
      return 'depart'
    case 'on_the_way':
      return 'arrive'
    case 'arrived':
      return 'start'
    case 'in_service':
      return 'complete'
    default:
      return null
  }
}

export function getOrderActionLabel(action: ProviderOrderAction): string {
  switch (action) {
    case 'accept':
      return '确认接单'
    case 'depart':
      return '立即出发'
    case 'arrive':
      return '确认到达'
    case 'start':
      return '开始服务'
    case 'complete':
      return '服务完单'
    default:
      return '执行操作'
  }
}
