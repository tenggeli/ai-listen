export interface MockOrderRecord {
  id: string
  providerId: string
  providerName: string
  serviceItemId: string
  serviceItemTitle: string
  amount: number
  currency: string
  status: 'paid'
  createdAt: string
}

const STORAGE_KEY = 'listen_user_web_mock_orders'

export function createPaidOrder(input: Omit<MockOrderRecord, 'id' | 'status' | 'createdAt'>): MockOrderRecord {
  const next: MockOrderRecord = {
    ...input,
    id: buildOrderId(),
    status: 'paid',
    createdAt: new Date().toISOString()
  }
  const all = readOrders()
  all.unshift(next)
  window.sessionStorage.setItem(STORAGE_KEY, JSON.stringify(all))
  return next
}

export function getOrderById(orderId: string): MockOrderRecord | null {
  return readOrders().find((item) => item.id === orderId) ?? null
}

function readOrders(): MockOrderRecord[] {
  const raw = window.sessionStorage.getItem(STORAGE_KEY)
  if (!raw) {
    return []
  }
  try {
    const parsed = JSON.parse(raw) as MockOrderRecord[]
    if (!Array.isArray(parsed)) {
      return []
    }
    return parsed
  } catch {
    return []
  }
}

function buildOrderId(): string {
  const seed = Math.random().toString(36).slice(2, 10)
  return `ord_${Date.now()}_${seed}`
}
