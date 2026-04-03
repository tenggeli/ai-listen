import { describe, expect, it } from 'vitest'
import { ApiError } from '../../api/ApiError'
import { toOrderDetailErrorMessage } from './orderError'

describe('toOrderDetailErrorMessage', () => {
  it('returns user-friendly message for 404 detail', () => {
    expect(toOrderDetailErrorMessage(new ApiError('not found', 404))).toBe('订单不存在或无权查看。')
  })

  it('keeps backend message for other errors', () => {
    expect(toOrderDetailErrorMessage(new ApiError('网络请求超时，请稍后重试', 408))).toBe('网络请求超时，请稍后重试')
  })
})
