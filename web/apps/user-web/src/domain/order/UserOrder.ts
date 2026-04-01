export type UserOrderStatus = 'created' | 'paid'

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
