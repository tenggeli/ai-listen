import { reactive } from 'vue'
import type { ProviderAdminApi } from '../../api/ProviderAdminApi'
import { PageLoadState } from '../../domain/provider/PageLoadState'
import { ProviderReviewStatus } from '../../domain/provider/ProviderReviewStatus'
import type { ProviderSummary } from '../../domain/provider/ProviderSummary'
import type { ProviderDetail } from '../../domain/provider/ProviderDetail'

export interface ProviderReviewState {
  listState: PageLoadState
  detailState: PageLoadState
  actionState: PageLoadState
  errorMessage: string
  activeFilter: string
  providers: ProviderSummary[]
  selectedProviderId: string
  selectedProviderDetail: ProviderDetail | null
}

export class ProviderReviewViewModel {
  readonly state: ProviderReviewState = reactive({
    listState: PageLoadState.Idle,
    detailState: PageLoadState.Idle,
    actionState: PageLoadState.Idle,
    errorMessage: '',
    activeFilter: '',
    providers: [],
    selectedProviderId: '',
    selectedProviderDetail: null
  })

  constructor(private readonly api: ProviderAdminApi) {}

  async initialize(): Promise<void> {
    await this.loadList()
  }

  async loadList(): Promise<void> {
    this.state.listState = PageLoadState.Loading
    this.state.errorMessage = ''
    try {
      const result = await this.api.listProviders(this.state.activeFilter)
      this.state.providers = result.items
      this.state.listState = result.items.length === 0 ? PageLoadState.Empty : PageLoadState.Success

      if (result.items.length > 0) {
        this.state.selectedProviderId = result.items[0].id
        await this.loadDetail(result.items[0].id)
      } else {
        this.state.selectedProviderId = ''
        this.state.selectedProviderDetail = null
        this.state.detailState = PageLoadState.Empty
      }
    } catch (error) {
      this.state.listState = PageLoadState.Error
      this.state.errorMessage = error instanceof Error ? error.message : '加载服务方失败'
    }
  }

  async changeFilter(status: string): Promise<void> {
    this.state.activeFilter = status
    await this.loadList()
  }

  async selectProvider(providerId: string): Promise<void> {
    this.state.selectedProviderId = providerId
    await this.loadDetail(providerId)
  }

  async review(action: 'approve' | 'reject' | 'require-supplement'): Promise<void> {
    if (!this.state.selectedProviderId) {
      return
    }
    this.state.actionState = PageLoadState.Loading
    this.state.errorMessage = ''
    try {
      const nextStatus = await this.api.review(this.state.selectedProviderId, action, '')
      this.patchListStatus(this.state.selectedProviderId, nextStatus)
      if (this.state.selectedProviderDetail) {
        this.state.selectedProviderDetail = {
          ...this.state.selectedProviderDetail,
          reviewStatus: nextStatus
        }
      }
      this.state.actionState = PageLoadState.Success
    } catch (error) {
      this.state.actionState = PageLoadState.Error
      this.state.errorMessage = error instanceof Error ? error.message : '审核操作失败'
    }
  }

  private async loadDetail(providerId: string): Promise<void> {
    this.state.detailState = PageLoadState.Loading
    try {
      this.state.selectedProviderDetail = await this.api.getProviderDetail(providerId)
      this.state.detailState = PageLoadState.Success
    } catch (error) {
      this.state.detailState = PageLoadState.Error
      this.state.errorMessage = error instanceof Error ? error.message : '加载详情失败'
    }
  }

  private patchListStatus(providerId: string, status: ProviderReviewStatus): void {
    this.state.providers = this.state.providers.map((item) =>
      item.id === providerId
        ? { ...item, reviewStatus: status }
        : item
    )
  }
}
