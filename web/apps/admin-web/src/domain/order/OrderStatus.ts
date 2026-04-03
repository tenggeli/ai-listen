export type AdminOrderStatus =
  | 'created'
  | 'paid'
  | 'accepted'
  | 'on_the_way'
  | 'arrived'
  | 'in_service'
  | 'completed'
  | 'after_sale_processing'
  | 'closed'

export const ADMIN_ORDER_STATUSES: AdminOrderStatus[] = [
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

const ADMIN_ORDER_STATUS_REASON_MAP: Record<AdminOrderStatus, string> = {
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

export function isAdminOrderStatus(value: string): value is AdminOrderStatus {
  return ADMIN_ORDER_STATUSES.includes(value as AdminOrderStatus)
}

export function getAdminOrderStatusReason(status: AdminOrderStatus): string {
  return ADMIN_ORDER_STATUS_REASON_MAP[status]
}
