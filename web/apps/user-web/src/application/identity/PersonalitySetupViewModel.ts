import { reactive } from 'vue'
import type { AuthApi } from '../../api/AuthApi'
import { loadSession, updateSessionPatch } from './AuthSession'

export interface PersonalitySetupState {
  loading: boolean
  saving: boolean
  selectedMbti: string
  selectedTags: string[]
  errorMessage: string
}

export class PersonalitySetupViewModel {
  readonly state: PersonalitySetupState = reactive({
    loading: false,
    saving: false,
    selectedMbti: '',
    selectedTags: [],
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
      this.state.selectedMbti = me.mbti
      this.state.selectedTags = [...me.interestTags]
    } finally {
      this.state.loading = false
    }
  }

  toggleTag(tag: string): void {
    const exists = this.state.selectedTags.includes(tag)
    if (exists) {
      this.state.selectedTags = this.state.selectedTags.filter((item) => item !== tag)
      return
    }
    this.state.selectedTags = [...this.state.selectedTags, tag]
  }

  toggleMbti(mbti: string): void {
    this.state.selectedMbti = this.state.selectedMbti === mbti ? '' : mbti
  }

  async submitSave(): Promise<string> {
    const session = loadSession()
    if (!session) {
      throw new Error('会话已过期，请重新登录')
    }

    if (!this.state.selectedMbti && this.state.selectedTags.length === 0) {
      throw new Error('请至少选择 MBTI 或兴趣标签')
    }

    this.state.saving = true
    this.state.errorMessage = ''
    try {
      const me = await this.api.savePersonality(session.accessToken, {
        mbti: this.state.selectedMbti,
        interestTags: this.state.selectedTags
      })
      updateSessionPatch({ personalityCompleted: me.personalityCompleted })
      return '/home'
    } finally {
      this.state.saving = false
    }
  }

  async submitSkip(): Promise<string> {
    const session = loadSession()
    if (!session) {
      throw new Error('会话已过期，请重新登录')
    }

    this.state.saving = true
    this.state.errorMessage = ''
    try {
      const me = await this.api.skipPersonality(session.accessToken)
      updateSessionPatch({ personalityCompleted: me.personalityCompleted })
      return '/home'
    } finally {
      this.state.saving = false
    }
  }
}

