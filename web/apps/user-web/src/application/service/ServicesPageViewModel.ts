import { reactive } from 'vue'
import type { ListPublicProvidersOutput, ServiceDiscoveryApi } from '../../api/ServiceDiscoveryApi'
import { PageLoadState } from '../../domain/ai/PageLoadState'
import type { ProviderPublicProfile } from '../../domain/service/ProviderPublicProfile'
import type { ServiceCategory } from '../../domain/service/ServiceCategory'

export interface ServicesPageState {
  pageState: PageLoadState
  categories: ServiceCategory[]
  selectedCategoryId: string
  searchKeyword: string
  providers: ProviderPublicProfile[]
  total: number
  errorMessage: string
}

export class ServicesPageViewModel {
  readonly state: ServicesPageState = reactive({
    pageState: PageLoadState.Idle,
    categories: [],
    selectedCategoryId: 'cat_all',
    searchKeyword: '',
    providers: [],
    total: 0,
    errorMessage: ''
  })

  constructor(private readonly api: ServiceDiscoveryApi) {}

  async initialize(): Promise<void> {
    this.state.pageState = PageLoadState.Loading
    this.state.errorMessage = ''

    try {
      const categories = await this.api.getCategories()
      this.state.categories = categories
      if (categories.length > 0 && !categories.some((item) => item.id === this.state.selectedCategoryId)) {
        this.state.selectedCategoryId = categories[0].id
      }
      const providers = await this.loadProviders()
      this.applyProviders(providers)
    } catch (error) {
      this.state.pageState = PageLoadState.Error
      this.state.errorMessage = error instanceof Error ? error.message : '加载服务页失败'
    }
  }

  async selectCategory(categoryId: string): Promise<void> {
    this.state.selectedCategoryId = categoryId
    await this.search()
  }

  async search(): Promise<void> {
    this.state.pageState = PageLoadState.Loading
    this.state.errorMessage = ''

    try {
      const providers = await this.loadProviders()
      this.applyProviders(providers)
    } catch (error) {
      this.state.pageState = PageLoadState.Error
      this.state.errorMessage = error instanceof Error ? error.message : '加载服务方失败'
    }
  }

  private async loadProviders(): Promise<ListPublicProvidersOutput> {
    return this.api.listProviders({
      categoryId: this.state.selectedCategoryId,
      keyword: this.state.searchKeyword,
      page: 1,
      pageSize: 20
    })
  }

  private applyProviders(output: ListPublicProvidersOutput): void {
    this.state.providers = output.items
    this.state.total = output.total
    this.state.pageState = output.items.length > 0 ? PageLoadState.Success : PageLoadState.Empty
  }
}
