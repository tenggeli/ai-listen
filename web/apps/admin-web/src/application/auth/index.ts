import { HttpAdminAuthApi } from '../../api/AdminAuthApi'
import { AuthService } from './AuthService'
import { AuthSessionStore } from './AuthSessionStore'

const authApi = new HttpAdminAuthApi('/api/v1/admin/auth')
export const authSessionStore = new AuthSessionStore()
export const authService = new AuthService(authApi, authSessionStore)
