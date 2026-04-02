import { HttpProviderAuthApi } from '../../api/ProviderAuthApi'
import { AuthService } from './AuthService'
import { AuthSessionStore } from './AuthSessionStore'

const authApi = new HttpProviderAuthApi('/api/v1/provider')
export const authSessionStore = new AuthSessionStore()
export const authService = new AuthService(authApi, authSessionStore)
