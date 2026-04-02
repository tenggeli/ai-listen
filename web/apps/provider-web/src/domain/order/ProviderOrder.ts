export type ProviderOrderStatus = 'created' | 'paid' | 'accepted' | 'on_the_way' | 'arrived' | 'in_service' | 'completed'

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
  createdAt: string
  paidAt: string | null
}
