export type ProviderOrderStatus =
  | 'created'
  | 'paid'
  | 'accepted'
  | 'on_the_way'
  | 'arrived'
  | 'in_service'
  | 'completed'
  | 'after_sale_processing'
  | 'closed'

export const PROVIDER_ORDER_STATUSES: ProviderOrderStatus[] = [
  'created',
  'paid',
  'accepted',
  'on_the_way',
  'arrived',
  'in_service',
  'completed',
  'after_sale_processing',
  'closed'
]

const PROVIDER_ORDER_STATUS_REASON_MAP: Record<ProviderOrderStatus, string> = {
  created: '待支付',
  paid: '待服务方接单',
  accepted: '服务方已接单',
  on_the_way: '服务方出发中',
  arrived: '服务方已到达，待开始服务',
  in_service: '服务进行中',
  completed: '服务已完成',
  after_sale_processing: '订单售后处理中',
  closed: '订单已关闭'
}

export function isProviderOrderStatus(value: string): value is ProviderOrderStatus {
  return PROVIDER_ORDER_STATUSES.includes(value as ProviderOrderStatus)
}

export function getProviderOrderStatusReason(status: ProviderOrderStatus): string {
  return PROVIDER_ORDER_STATUS_REASON_MAP[status]
}

export interface ProviderOrder {
  id: string
  userId: string
  providerId: string
  providerName: string
  serviceItemId: string
  serviceItemTitle: string
  amount: number
  currency: string
  status: ProviderOrderStatus
  statusReason: string
  statusActionReason: string
  statusUpdatedAt: string | null
  createdAt: string
  paidAt: string | null
}
