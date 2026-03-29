import { reactive } from 'vue'
import type { ServiceDiscoveryApi } from '../../api/ServiceDiscoveryApi'
import { PageLoadState } from '../../domain/ai/PageLoadState'
import type { ProviderPublicProfile } from '../../domain/service/ProviderPublicProfile'
import type { ServiceItem } from '../../domain/service/ServiceItem'

export interface ProviderDetailPageState {
  pageState: PageLoadState
  provider: ProviderPublicProfile | null
  serviceItems: ServiceItem[]
  selectedServiceItemId: string
  errorMessage: string
}

export class ProviderDetailPageViewModel {
  readonly state: ProviderDetailPageState = reactive({
    pageState: PageLoadState.Idle,
    provider: null,
    serviceItems: [],
    selectedServiceItemId: '',
    errorMessage: ''
  })

  constructor(private readonly api: ServiceDiscoveryApi) {}

  async initialize(providerId: string): Promise<void> {
    this.state.pageState = PageLoadState.Loading
    this.state.errorMessage = ''

    try {
      const [provider, serviceItems] = await Promise.all([
        this.api.getProvider(providerId),
        this.api.getProviderServiceItems(providerId)
      ])
      this.state.provider = provider
      this.state.serviceItems = serviceItems
      this.state.selectedServiceItemId = serviceItems[0]?.id ?? ''
      this.state.pageState = provider && serviceItems.length > 0 ? PageLoadState.Success : PageLoadState.Empty
    } catch (error) {
      this.state.pageState = PageLoadState.Error
      this.state.errorMessage = error instanceof Error ? error.message : '加载服务方详情失败'
    }
  }

  selectServiceItem(serviceItemId: string): void {
    this.state.selectedServiceItemId = serviceItemId
  }
}
