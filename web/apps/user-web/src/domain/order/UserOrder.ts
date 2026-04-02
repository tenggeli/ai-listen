export type UserOrderStatus =
  | 'created'
  | 'paid'
  | 'accepted'
  | 'on_the_way'
  | 'arrived'
  | 'in_service'
  | 'completed'
  | 'after_sale_processing'
  | 'closed'

export interface UserOrder {
  id: string
  userId: string
  providerId: string
  providerName: string
  serviceItemId: string
  serviceItemTitle: string
  amount: number
  currency: string
  status: UserOrderStatus
  createdAt: string
  paidAt: string | null
}
