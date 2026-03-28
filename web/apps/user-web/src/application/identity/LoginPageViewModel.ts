import { reactive } from 'vue'
import type { AuthApi } from '../../api/AuthApi'
import { nextOnboardingRoute, saveSession } from './AuthSession'

type LoginMode = 'wechat' | 'sms'

export interface LoginPageState {
  mode: LoginMode
  phone: string
  verifyCode: string
  agreementAccepted: boolean
  submitting: boolean
  errorMessage: string
}

export class LoginPageViewModel {
  readonly state: LoginPageState = reactive({
    mode: 'wechat',
    phone: '13800113800',
    verifyCode: '123',
    agreementAccepted: false,
    submitting: false,
    errorMessage: ''
  })

  constructor(private readonly api: AuthApi) {}

  setMode(mode: LoginMode): void {
    this.state.mode = mode
    this.state.errorMessage = ''
  }

  toggleAgreement(): void {
    this.state.agreementAccepted = !this.state.agreementAccepted
    this.state.errorMessage = ''
  }

  async submitSmsLogin(): Promise<string> {
    if (!this.state.agreementAccepted) {
      throw new Error('请先勾选用户协议')
    }
    if (!/^1\d{10}$/.test(this.state.phone.replace(/\s/g, ''))) {
      throw new Error('手机号格式错误')
    }
    if (!this.state.verifyCode.trim()) {
      throw new Error('请输入验证码')
    }

    this.state.submitting = true
    this.state.errorMessage = ''
    try {
      const identity = await this.api.loginBySms(this.state.phone.replace(/\s/g, ''), this.state.verifyCode.trim(), true)
      const me = await this.api.getMe(identity.accessToken)
      const session = saveSession(identity, me)
      return nextOnboardingRoute(session)
    } finally {
      this.state.submitting = false
    }
  }

  async submitWechatLogin(): Promise<string> {
    if (!this.state.agreementAccepted) {
      throw new Error('请先勾选用户协议')
    }

    this.state.submitting = true
    this.state.errorMessage = ''
    try {
      const identity = await this.api.loginByWechatMock(`wx_${Date.now()}`, true)
      const me = await this.api.getMe(identity.accessToken)
      const session = saveSession(identity, me)
      return nextOnboardingRoute(session)
    } finally {
      this.state.submitting = false
    }
  }
}
