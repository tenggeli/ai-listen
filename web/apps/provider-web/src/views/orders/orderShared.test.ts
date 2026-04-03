import { describe, expect, it } from 'vitest'
import type { ProviderOrder } from '../../domain/order/ProviderOrder'
import { filterOrdersByStatus, getTotalPages } from './orderShared'

const sampleOrders: ProviderOrder[] = [
  {
    id: 'ord-1',
    userId: 'u1',
    providerId: 'p1',
    providerName: 'Provider',
    serviceItemId: 'svc-1',
    serviceItemTitle: '心理疏导',
    amount: 200,
    currency: 'CNY',
    status: 'paid',
    createdAt: '2026-04-01T10:00:00Z',
    paidAt: '2026-04-01T10:05:00Z'
  },
  {
    id: 'ord-2',
    userId: 'u2',
    providerId: 'p1',
    providerName: 'Provider',
    serviceItemId: 'svc-2',
    serviceItemTitle: '陪伴',
    amount: 180,
    currency: 'CNY',
    status: 'completed',
    createdAt: '2026-04-01T11:00:00Z',
    paidAt: '2026-04-01T11:05:00Z'
  }
]

describe('order shared helpers', () => {
  it('supports empty state filtering', () => {
    expect(filterOrdersByStatus([], 'all')).toEqual([])
    expect(filterOrdersByStatus([], 'paid')).toEqual([])
  })

  it('filters list by status', () => {
    const paidOnly = filterOrdersByStatus(sampleOrders, 'paid')
    expect(paidOnly).toHaveLength(1)
    expect(paidOnly[0].id).toBe('ord-1')
  })

  it('handles pagination boundaries', () => {
    expect(getTotalPages(0, 10)).toBe(1)
    expect(getTotalPages(1, 10)).toBe(1)
    expect(getTotalPages(10, 10)).toBe(1)
    expect(getTotalPages(11, 10)).toBe(2)
  })
})
