package store

type Store interface {
	IssueSMSCode(mobile string) string
	LoginBySMS(mobile, code string) (*User, string, string, error)
	RefreshToken(refreshToken string) (string, error)
	AdminLogin(username, password string) (*AdminUser, string, error)
	AdminByToken(raw string) (*AdminUser, error)
	AdminUserByID(adminUserID uint64) (*AdminUser, error)
	AdminPermissionsByRoles(roleKeys []string) []string
	AdminHasPermission(adminUserID uint64, permission string) bool
	AdminRoles() []*AdminRole
	AdminPermissions() []*AdminPermission
	UpdateAdminUserRoles(adminUserID uint64, roleKeys []string) (*AdminUser, error)
	UserByToken(raw string) (*User, error)
	GetUser(userID uint64) (*User, error)
	UpdateUser(userID uint64, updater func(*User)) (*User, error)
	ServiceItems() []*ServiceItem
	GetServiceItem(id uint64) (*ServiceItem, error)
	ApplyProvider(userID uint64, realName, idCardNo string) (*Provider, error)
	ProviderByUserID(userID uint64) (*Provider, error)
	ProviderByID(providerID uint64) (*Provider, error)
	Providers() []*Provider
	UpdateProviderProfile(userID uint64, displayName, intro string, tags []string) (*Provider, error)
	UpdateProviderServiceItems(userID uint64, items []ProviderServiceItem) (*Provider, error)
	UpdateProviderWorkStatus(userID uint64, workStatus int) (*Provider, error)
	ApproveProvider(providerID uint64, remark string) (*Provider, error)
	RejectProvider(providerID uint64, remark string) (*Provider, error)
	CreateOrder(userID uint64, req CreateOrderInput) (*Order, error)
	OrdersByUser(userID uint64) []*Order
	OrdersByProvider(userID uint64) ([]*Order, error)
	GetOrder(orderID uint64) (*Order, error)
	CreatePayment(orderID uint64) (*Payment, *Order, error)
	CancelOrder(orderID, userID uint64, reason string) (*Order, error)
	ProviderAcceptOrder(userID, orderID uint64) (*Order, error)
	ProviderDepartOrder(userID, orderID uint64) (*Order, error)
	ProviderArriveOrder(userID, orderID uint64) (*Order, error)
	StartOrder(userID, orderID uint64, remark string) (*Order, error)
	ProviderFinishOrder(userID, orderID uint64) (*Order, error)
	ConfirmFinishOrder(userID, orderID uint64) (*Order, error)
}

var defaultStore Store = NewMemoryStore()

func Default() Store { return defaultStore }

func SetDefault(s Store) {
	if s == nil {
		return
	}
	defaultStore = s
}

func ResetDefaultForTest() { defaultStore = NewMemoryStore() }
