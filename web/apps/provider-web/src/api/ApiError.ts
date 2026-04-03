export class ApiError extends Error {
  constructor(
    message: string,
    readonly statusCode: number
  ) {
    super(message)
    this.name = 'ApiError'
  }
}

export function isUnauthorizedError(error: unknown): boolean {
  return error instanceof ApiError && error.statusCode === 401
}
