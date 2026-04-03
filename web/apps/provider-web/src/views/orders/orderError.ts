import { ApiError } from '../../api/ApiError'

export function toOrderDetailErrorMessage(error: unknown): string {
  if (error instanceof ApiError && error.statusCode === 404) {
    return '订单不存在或无权查看。'
  }
  return error instanceof Error ? error.message : '订单详情加载失败'
}
