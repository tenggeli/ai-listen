import { reactive } from 'vue'
import type { ServiceItemAdminApi } from '../../api/ServiceItemAdminApi'
import { PageLoadState } from '../../domain/provider/PageLoadState'
import type { ServiceItemDetail } from '../../domain/service_item/ServiceItemDetail'
import type { ServiceItemSummary } from '../../domain/service_item/ServiceItemSummary'

export interface ServiceItemManageState {
  listState: PageLoadState
  detailState: PageLoadState
  actionState: PageLoadState
  errorMessage: string
  providerId: string
  categoryId: string
  status: string
  keyword: string
  items: ServiceItemSummary[]
  selectedServiceItemId: string
  selectedServiceItemDetail: ServiceItemDetail | null
}

export class ServiceItemManageViewModel {
  readonly state: ServiceItemManageState = reactive({
    listState: PageLoadState.Idle,
    detailState: PageLoadState.Idle,
    actionState: PageLoadState.Idle,
    errorMessage: '',
    providerId: '',
    categoryId: '',
    status: '',
    keyword: '',
    items: [],
    selectedServiceItemId: '',
    selectedServiceItemDetail: null
  })

  constructor(private readonly api: ServiceItemAdminApi) {}

  async initialize(): Promise<void> {
    await this.loadList()
  }

  async loadList(): Promise<void> {
    this.state.listState = PageLoadState.Loading
    this.state.errorMessage = ''
    try {
      const result = await this.api.listServiceItems({
        providerId: this.state.providerId,
        categoryId: this.state.categoryId,
        status: this.state.status,
        keyword: this.state.keyword
      })
      this.state.items = result.items
      this.state.listState = result.items.length === 0 ? PageLoadState.Empty : PageLoadState.Success

      if (result.items.length > 0) {
        const selected = result.items.find((item) => item.id === this.state.selectedServiceItemId) ?? result.items[0]
        this.state.selectedServiceItemId = selected.id
        await this.loadDetail(selected.id)
      } else {
        this.state.selectedServiceItemId = ''
        this.state.selectedServiceItemDetail = null
        this.state.detailState = PageLoadState.Empty
      }
    } catch (error) {
      this.state.listState = PageLoadState.Error
      this.state.errorMessage = error instanceof Error ? error.message : '加载服务项目失败'
    }
  }

  async selectServiceItem(serviceItemId: string): Promise<void> {
    this.state.selectedServiceItemId = serviceItemId
    await this.loadDetail(serviceItemId)
  }

  async applyFilter(): Promise<void> {
    await this.loadList()
  }

  async updateStatus(action: 'activate' | 'deactivate'): Promise<void> {
    if (!this.state.selectedServiceItemId) {
      return
    }
    this.state.actionState = PageLoadState.Loading
    this.state.errorMessage = ''
    try {
      const status = await this.api.updateStatus(this.state.selectedServiceItemId, action)
      this.patchListStatus(this.state.selectedServiceItemId, status)
      if (this.state.selectedServiceItemDetail) {
        this.state.selectedServiceItemDetail = {
          ...this.state.selectedServiceItemDetail,
          serviceStatus: status
        }
      }
      this.state.actionState = PageLoadState.Success
    } catch (error) {
      this.state.actionState = PageLoadState.Error
      this.state.errorMessage = error instanceof Error ? error.message : '更新状态失败'
    }
  }

  private async loadDetail(serviceItemId: string): Promise<void> {
    this.state.detailState = PageLoadState.Loading
    this.state.errorMessage = ''
    try {
      this.state.selectedServiceItemDetail = await this.api.getServiceItemDetail(serviceItemId)
      this.state.detailState = PageLoadState.Success
    } catch (error) {
      this.state.detailState = PageLoadState.Error
      this.state.errorMessage = error instanceof Error ? error.message : '加载详情失败'
    }
  }

  private patchListStatus(serviceItemId: string, status: string): void {
    this.state.items = this.state.items.map((item) =>
      item.id === serviceItemId ? { ...item, serviceStatus: status } : item
    )
  }
}

