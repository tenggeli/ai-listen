import { reactive } from 'vue'
import type { AuthApi } from '../../api/AuthApi'
import { loadSession, updateSessionPatch } from './AuthSession'

export interface ProfileSetupState {
  loading: boolean
  saving: boolean
  nickname: string
  avatarUrl: string
  gender: string
  ageRange: string
  city: string
  bio: string
  errorMessage: string
}

const SKIP_PROFILE_DEFAULTS = {
  nickname: 'listener',
  gender: 'unknown',
  ageRange: 'unknown',
  city: 'unknown',
  bio: '',
  avatarUrl: ''
}

export class ProfileSetupViewModel {
  readonly state: ProfileSetupState = reactive({
    loading: false,
    saving: false,
    nickname: '',
    avatarUrl: '',
    gender: '',
    ageRange: '',
    city: '',
    bio: '',
    errorMessage: ''
  })

  constructor(private readonly api: AuthApi) {}

  async initialize(): Promise<void> {
    const session = loadSession()
    if (!session) {
      throw new Error('请先登录')
    }
    this.state.loading = true
    this.state.errorMessage = ''
    try {
      const me = await this.api.getMe(session.accessToken)
      this.state.nickname = me.nickname
      this.state.avatarUrl = me.avatarUrl
      this.state.gender = me.gender
      this.state.ageRange = me.ageRange
      this.state.city = me.city
      this.state.bio = me.bio
    } finally {
      this.state.loading = false
    }
  }

  async submitSave(): Promise<string> {
    if (!this.state.nickname.trim()) {
      throw new Error('请填写昵称')
    }
    if (!this.state.gender.trim()) {
      throw new Error('请选择性别')
    }
    if (!this.state.ageRange.trim()) {
      throw new Error('请选择年龄段')
    }
    if (!this.state.city.trim()) {
      throw new Error('请填写所在城市')
    }
    return this.saveInternal(false)
  }

  async submitSkip(): Promise<string> {
    this.state.nickname = this.state.nickname.trim() || SKIP_PROFILE_DEFAULTS.nickname
    this.state.gender = this.state.gender.trim() || SKIP_PROFILE_DEFAULTS.gender
    this.state.ageRange = this.state.ageRange.trim() || SKIP_PROFILE_DEFAULTS.ageRange
    this.state.city = this.state.city.trim() || SKIP_PROFILE_DEFAULTS.city
    if (!this.state.bio.trim()) {
      this.state.bio = SKIP_PROFILE_DEFAULTS.bio
    }
    return this.saveInternal(true)
  }

  private async saveInternal(skipMode: boolean): Promise<string> {
    const session = loadSession()
    if (!session) {
      throw new Error('会话已过期，请重新登录')
    }

    this.state.saving = true
    this.state.errorMessage = ''
    try {
      try {
        const me = await this.api.saveProfile(session.accessToken, {
          nickname: this.state.nickname.trim(),
          avatarUrl: this.state.avatarUrl.trim(),
          gender: this.state.gender.trim(),
          ageRange: this.state.ageRange.trim(),
          city: this.state.city.trim(),
          bio: this.state.bio.trim(),
          genderChangeConfirmed: false
        })
        updateSessionPatch({ profileCompleted: me.profileCompleted })
      } catch (error) {
        const message = error instanceof Error ? error.message : '保存资料失败'
        if (!skipMode && message.includes('gender change requires confirmation')) {
          const confirmed = window.confirm('检测到你在修改性别，确认继续修改吗？')
          if (!confirmed) {
            throw new Error('已取消性别修改')
          }
          const me = await this.api.saveProfile(session.accessToken, {
            nickname: this.state.nickname.trim(),
            avatarUrl: this.state.avatarUrl.trim(),
            gender: this.state.gender.trim(),
            ageRange: this.state.ageRange.trim(),
            city: this.state.city.trim(),
            bio: this.state.bio.trim(),
            genderChangeConfirmed: true
          })
          updateSessionPatch({ profileCompleted: me.profileCompleted })
        } else {
          throw error
        }
      }
      return '/personality/setup'
    } finally {
      this.state.saving = false
    }
  }
}

