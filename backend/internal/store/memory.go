package store

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

const (
	OrderStatusPendingPayment = 10
	OrderStatusPendingAccept  = 20
	OrderStatusAccepted       = 30
	OrderStatusDeparted       = 40
	OrderStatusArrived        = 50
	OrderStatusServing        = 60
	OrderStatusPendingFinish  = 70
	OrderStatusCompleted      = 80
	OrderStatusCanceled       = 90
)

var (
	ErrUnauthorized          = errors.New("unauthorized")
	ErrUserNotFound          = errors.New("user not found")
	ErrAdminNotFound         = errors.New("admin user not found")
	ErrProviderNotFound      = errors.New("provider not found")
	ErrRoleNotFound          = errors.New("admin role not found")
	ErrProviderAuditRequired = errors.New("provider audit not passed")
	ErrOrderNotFound         = errors.New("order not found")
	ErrInvalidOrderStatus    = errors.New("invalid order status")
	ErrInvalidSMSCode        = errors.New("invalid sms code")
	ErrServiceItemNotFound   = errors.New("service item not found")
)

type User struct {
	ID        uint64    `json:"id"`
	Mobile    string    `json:"mobile"`
	Nickname  string    `json:"nickname"`
	Avatar    string    `json:"avatar"`
	Gender    int       `json:"gender"`
	Birthday  string    `json:"birthday"`
	CityCode  string    `json:"cityCode"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ProviderServiceItem struct {
	ServiceItemID uint64 `json:"serviceItemId"`
	PriceAmount   int64  `json:"priceAmount"`
	PriceUnit     string `json:"priceUnit"`
}

type Provider struct {
	ID           uint64                `json:"id"`
	UserID       uint64                `json:"userId"`
	ProviderNo   string                `json:"providerNo"`
	RealName     string                `json:"realName"`
	IDCardNo     string                `json:"idCardNo"`
	AuditStatus  int                   `json:"auditStatus"`
	AuditRemark  string                `json:"auditRemark"`
	WorkStatus   int                   `json:"workStatus"`
	DisplayName  string                `json:"displayName"`
	Intro        string                `json:"intro"`
	Tags         []string              `json:"tags"`
	ServiceItems []ProviderServiceItem `json:"serviceItems"`
	CreatedAt    time.Time             `json:"createdAt"`
	UpdatedAt    time.Time             `json:"updatedAt"`
}

type ServiceItem struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Unit     string `json:"unit"`
	MinPrice int64  `json:"minPrice"`
	MaxPrice int64  `json:"maxPrice"`
	Status   int    `json:"status"`
}

type OrderLog struct {
	FromStatus   int       `json:"fromStatus"`
	ToStatus     int       `json:"toStatus"`
	OperatorRole string    `json:"operatorRole"`
	OperatorID   uint64    `json:"operatorId"`
	Remark       string    `json:"remark"`
	CreatedAt    time.Time `json:"createdAt"`
}

type Order struct {
	ID              uint64     `json:"id"`
	OrderNo         string     `json:"orderNo"`
	UserID          uint64     `json:"userId"`
	ProviderID      uint64     `json:"providerId"`
	ServiceItemID   uint64     `json:"serviceItemId"`
	SceneText       string     `json:"sceneText"`
	CityCode        string     `json:"cityCode"`
	AddressText     string     `json:"addressText"`
	PlannedStartAt  string     `json:"plannedStartAt"`
	PlannedDuration int        `json:"plannedDuration"`
	Status          int        `json:"status"`
	PayAmount       int64      `json:"payAmount"`
	CancelReason    string     `json:"cancelReason,omitempty"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
	AcceptedAt      *time.Time `json:"acceptedAt,omitempty"`
	ArrivedAt       *time.Time `json:"arrivedAt,omitempty"`
	StartedAt       *time.Time `json:"startedAt,omitempty"`
	FinishedAt      *time.Time `json:"finishedAt,omitempty"`
	Logs            []OrderLog `json:"logs"`
}

type Payment struct {
	ID        uint64    `json:"id"`
	OrderID   uint64    `json:"orderId"`
	PayAmount int64     `json:"payAmount"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	PaidAt    time.Time `json:"paidAt"`
}

type CreateOrderInput struct {
	ProviderID      uint64
	ServiceItemID   uint64
	SceneText       string
	CityCode        string
	AddressText     string
	PlannedStartAt  string
	PlannedDuration int
}

type MemoryStore struct {
	mu sync.RWMutex

	nextUserID     uint64
	nextAdminID    uint64
	nextProviderID uint64
	nextOrderID    uint64
	nextPaymentID  uint64

	users          map[uint64]*User
	usersByMobile  map[string]*User
	providers      map[uint64]*Provider
	providerByUID  map[uint64]*Provider
	orders         map[uint64]*Order
	payments       map[uint64]*Payment
	smsCodes       map[string]string
	tokens         map[string]uint64
	refreshTokens  map[string]uint64
	adminUsers     map[uint64]*AdminUser
	adminByName    map[string]*AdminUser
	adminPasswords map[uint64]string
	adminTokens    map[string]uint64
	adminRoles     map[string]*AdminRole
	adminPerms     map[string]*AdminPermission
	rolePerms      map[string]map[string]struct{}
	serviceItems   map[uint64]*ServiceItem
}

func NewMemoryStore() *MemoryStore {
	s := &MemoryStore{
		nextUserID:     1,
		nextAdminID:    1,
		nextProviderID: 1,
		nextOrderID:    1,
		nextPaymentID:  1,
		users:          map[uint64]*User{},
		usersByMobile:  map[string]*User{},
		providers:      map[uint64]*Provider{},
		providerByUID:  map[uint64]*Provider{},
		orders:         map[uint64]*Order{},
		payments:       map[uint64]*Payment{},
		smsCodes:       map[string]string{},
		tokens:         map[string]uint64{},
		refreshTokens:  map[string]uint64{},
		adminUsers:     map[uint64]*AdminUser{},
		adminByName:    map[string]*AdminUser{},
		adminPasswords: map[uint64]string{},
		adminTokens:    map[string]uint64{},
		adminRoles:     map[string]*AdminRole{},
		adminPerms:     map[string]*AdminPermission{},
		rolePerms:      map[string]map[string]struct{}{},
		serviceItems:   map[uint64]*ServiceItem{},
	}

	s.serviceItems[1] = &ServiceItem{ID: 1, Name: "陪吃饭", Category: "陪伴", Unit: "小时", MinPrice: 10000, MaxPrice: 50000, Status: 1}
	s.serviceItems[2] = &ServiceItem{ID: 2, Name: "观影搭子", Category: "娱乐", Unit: "小时", MinPrice: 8000, MaxPrice: 40000, Status: 1}
	s.serviceItems[3] = &ServiceItem{ID: 3, Name: "心理疏导", Category: "情绪支持", Unit: "小时", MinPrice: 20000, MaxPrice: 100000, Status: 1}
	s.seedAdminRBAC()
	defaultAdmin := &AdminUser{
		ID:        s.nextAdminID,
		Username:  "admin",
		Nickname:  "系统管理员",
		Roles:     []string{"super_admin"},
		Status:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	s.nextAdminID++
	s.adminUsers[defaultAdmin.ID] = defaultAdmin
	s.adminByName[defaultAdmin.Username] = defaultAdmin
	s.adminPasswords[defaultAdmin.ID] = "admin123456"

	contentAdmin := &AdminUser{
		ID:        s.nextAdminID,
		Username:  "content_admin",
		Nickname:  "内容管理员",
		Roles:     []string{"content_admin"},
		Status:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	s.nextAdminID++
	s.adminUsers[contentAdmin.ID] = contentAdmin
	s.adminByName[contentAdmin.Username] = contentAdmin
	s.adminPasswords[contentAdmin.ID] = "admin123456"
	return s
}

func (s *MemoryStore) seedAdminRBAC() {
	now := time.Now()
	permissionSeed := []struct {
		Name          string
		PermissionKey string
		Category      string
	}{
		{Name: "看板总览", PermissionKey: "admin:dashboard:overview", Category: "dashboard"},
		{Name: "用户列表", PermissionKey: "admin:user:list", Category: "user"},
		{Name: "用户详情", PermissionKey: "admin:user:detail", Category: "user"},
		{Name: "用户状态更新", PermissionKey: "admin:user:status:update", Category: "user"},
		{Name: "服务方列表", PermissionKey: "admin:provider:list", Category: "provider"},
		{Name: "服务方详情", PermissionKey: "admin:provider:detail", Category: "provider"},
		{Name: "服务方审核通过", PermissionKey: "admin:provider:approve", Category: "provider"},
		{Name: "服务方审核拒绝", PermissionKey: "admin:provider:reject", Category: "provider"},
		{Name: "服务方状态更新", PermissionKey: "admin:provider:status:update", Category: "provider"},
		{Name: "服务项目列表", PermissionKey: "admin:service_item:list", Category: "service_item"},
		{Name: "服务项目创建", PermissionKey: "admin:service_item:create", Category: "service_item"},
		{Name: "服务项目更新", PermissionKey: "admin:service_item:update", Category: "service_item"},
		{Name: "服务项目删除", PermissionKey: "admin:service_item:delete", Category: "service_item"},
		{Name: "订单列表", PermissionKey: "admin:order:list", Category: "order"},
		{Name: "订单详情", PermissionKey: "admin:order:detail", Category: "order"},
		{Name: "订单人工完结", PermissionKey: "admin:order:manual_complete", Category: "order"},
		{Name: "订单退款", PermissionKey: "admin:order:refund", Category: "order"},
		{Name: "提现列表", PermissionKey: "admin:withdraw:list", Category: "withdraw"},
		{Name: "提现审核通过", PermissionKey: "admin:withdraw:approve", Category: "withdraw"},
		{Name: "提现审核拒绝", PermissionKey: "admin:withdraw:reject", Category: "withdraw"},
		{Name: "财务报表", PermissionKey: "admin:finance:reports", Category: "finance"},
		{Name: "帖子列表", PermissionKey: "admin:post:list", Category: "content"},
		{Name: "帖子隐藏", PermissionKey: "admin:post:hide", Category: "content"},
		{Name: "声音列表", PermissionKey: "admin:audio:list", Category: "content"},
		{Name: "声音下架", PermissionKey: "admin:audio:off_shelf", Category: "content"},
		{Name: "投诉列表", PermissionKey: "admin:complaint:list", Category: "complaint"},
		{Name: "投诉详情", PermissionKey: "admin:complaint:detail", Category: "complaint"},
		{Name: "投诉处理", PermissionKey: "admin:complaint:resolve", Category: "complaint"},
		{Name: "风控事件列表", PermissionKey: "admin:risk_event:list", Category: "risk"},
		{Name: "配置列表", PermissionKey: "admin:config:list", Category: "config"},
		{Name: "配置更新", PermissionKey: "admin:config:update", Category: "config"},
		{Name: "RBAC角色列表", PermissionKey: "admin:rbac:role:list", Category: "rbac"},
		{Name: "RBAC权限列表", PermissionKey: "admin:rbac:permission:list", Category: "rbac"},
		{Name: "管理员角色分配", PermissionKey: "admin:rbac:user_role:assign", Category: "rbac"},
	}
	for idx, item := range permissionSeed {
		s.adminPerms[item.PermissionKey] = &AdminPermission{
			ID:            uint64(idx + 1),
			Name:          item.Name,
			PermissionKey: item.PermissionKey,
			Category:      item.Category,
			Status:        1,
			CreatedAt:     now,
			UpdatedAt:     now,
		}
	}

	roleSeed := []struct {
		Name    string
		RoleKey string
	}{
		{Name: "超级管理员", RoleKey: "super_admin"},
		{Name: "运营管理员", RoleKey: "ops_admin"},
		{Name: "财务管理员", RoleKey: "finance_admin"},
		{Name: "内容管理员", RoleKey: "content_admin"},
	}
	for idx, item := range roleSeed {
		s.adminRoles[item.RoleKey] = &AdminRole{
			ID:        uint64(idx + 1),
			Name:      item.Name,
			RoleKey:   item.RoleKey,
			Status:    1,
			CreatedAt: now,
			UpdatedAt: now,
		}
		s.rolePerms[item.RoleKey] = map[string]struct{}{}
	}
	s.bindRolePermissions("super_admin", []string{"*"})
	s.bindRolePermissions("ops_admin", []string{
		"admin:dashboard:overview",
		"admin:user:list", "admin:user:detail", "admin:user:status:update",
		"admin:provider:list", "admin:provider:detail", "admin:provider:approve", "admin:provider:reject", "admin:provider:status:update",
		"admin:service_item:list", "admin:service_item:create", "admin:service_item:update", "admin:service_item:delete",
		"admin:order:list", "admin:order:detail", "admin:order:manual_complete",
		"admin:complaint:list", "admin:complaint:detail", "admin:complaint:resolve", "admin:risk_event:list",
	})
	s.bindRolePermissions("finance_admin", []string{
		"admin:dashboard:overview",
		"admin:order:list", "admin:order:detail", "admin:order:refund",
		"admin:withdraw:list", "admin:withdraw:approve", "admin:withdraw:reject", "admin:finance:reports",
	})
	s.bindRolePermissions("content_admin", []string{
		"admin:dashboard:overview",
		"admin:post:list", "admin:post:hide",
		"admin:audio:list", "admin:audio:off_shelf",
		"admin:complaint:list", "admin:complaint:detail", "admin:complaint:resolve",
	})
}

func (s *MemoryStore) bindRolePermissions(roleKey string, permissionKeys []string) {
	set, ok := s.rolePerms[roleKey]
	if !ok {
		set = map[string]struct{}{}
		s.rolePerms[roleKey] = set
	}
	for _, permissionKey := range permissionKeys {
		set[permissionKey] = struct{}{}
	}
}

func (s *MemoryStore) IssueSMSCode(mobile string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.smsCodes[mobile] = "123456"
	return "123456"
}

func (s *MemoryStore) LoginBySMS(mobile, code string) (*User, string, string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.smsCodes[mobile] != code {
		return nil, "", "", ErrInvalidSMSCode
	}

	user, ok := s.usersByMobile[mobile]
	if !ok {
		now := time.Now()
		user = &User{
			ID:        s.nextUserID,
			Mobile:    mobile,
			Nickname:  fmt.Sprintf("listen用户%04d", s.nextUserID),
			CityCode:  "310100",
			CreatedAt: now,
			UpdatedAt: now,
		}
		s.nextUserID++
		s.users[user.ID] = user
		s.usersByMobile[mobile] = user
	}

	token := fmt.Sprintf("token-%d", user.ID)
	refreshToken := fmt.Sprintf("refresh-%d", user.ID)
	s.tokens[token] = user.ID
	s.refreshTokens[refreshToken] = user.ID
	return cloneUser(user), token, refreshToken, nil
}

func (s *MemoryStore) RefreshToken(refreshToken string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	userID, ok := s.refreshTokens[refreshToken]
	if !ok {
		return "", ErrUnauthorized
	}
	token := fmt.Sprintf("token-%d-%d", userID, time.Now().UnixNano())
	s.tokens[token] = userID
	return token, nil
}

func (s *MemoryStore) AdminLogin(username, password string) (*AdminUser, string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	admin, ok := s.adminByName[username]
	if !ok {
		return nil, "", ErrUnauthorized
	}
	if s.adminPasswords[admin.ID] != password || admin.Status != 1 {
		return nil, "", ErrUnauthorized
	}

	token := fmt.Sprintf("admin-token-%d-%d", admin.ID, time.Now().UnixNano())
	s.adminTokens[token] = admin.ID
	return cloneAdmin(admin), token, nil
}

func (s *MemoryStore) AdminByToken(raw string) (*AdminUser, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	token := strings.TrimSpace(strings.TrimPrefix(raw, "Bearer "))
	adminID, ok := s.adminTokens[token]
	if !ok {
		return nil, ErrUnauthorized
	}
	admin, ok := s.adminUsers[adminID]
	if !ok {
		return nil, ErrUnauthorized
	}
	return cloneAdmin(admin), nil
}

func (s *MemoryStore) AdminUserByID(adminUserID uint64) (*AdminUser, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	admin, ok := s.adminUsers[adminUserID]
	if !ok {
		return nil, ErrAdminNotFound
	}
	return cloneAdmin(admin), nil
}

func (s *MemoryStore) AdminPermissionsByRoles(roleKeys []string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	permissionSet := map[string]struct{}{}
	for _, roleKey := range roleKeys {
		for permissionKey := range s.rolePerms[roleKey] {
			if permissionKey == "*" {
				for key := range s.adminPerms {
					permissionSet[key] = struct{}{}
				}
				continue
			}
			permissionSet[permissionKey] = struct{}{}
		}
	}
	result := make([]string, 0, len(permissionSet))
	for key := range permissionSet {
		result = append(result, key)
	}
	return result
}

func (s *MemoryStore) AdminHasPermission(adminUserID uint64, permission string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	admin, ok := s.adminUsers[adminUserID]
	if !ok {
		return false
	}
	for _, roleKey := range admin.Roles {
		if _, ok := s.rolePerms[roleKey]["*"]; ok {
			return true
		}
		if _, ok := s.rolePerms[roleKey][permission]; ok {
			return true
		}
	}
	return false
}

func (s *MemoryStore) AdminRoles() []*AdminRole {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*AdminRole, 0, len(s.adminRoles))
	for roleKey, role := range s.adminRoles {
		copyRole := cloneAdminRole(role)
		for permissionKey := range s.rolePerms[roleKey] {
			copyRole.PermissionKeys = append(copyRole.PermissionKeys, permissionKey)
		}
		result = append(result, copyRole)
	}
	return result
}

func (s *MemoryStore) AdminPermissions() []*AdminPermission {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*AdminPermission, 0, len(s.adminPerms))
	for _, permission := range s.adminPerms {
		copyPermission := *permission
		result = append(result, &copyPermission)
	}
	return result
}

func (s *MemoryStore) UpdateAdminUserRoles(adminUserID uint64, roleKeys []string) (*AdminUser, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	admin, ok := s.adminUsers[adminUserID]
	if !ok {
		return nil, ErrAdminNotFound
	}

	normalized := make([]string, 0, len(roleKeys))
	seen := map[string]struct{}{}
	for _, roleKey := range roleKeys {
		if _, exists := s.adminRoles[roleKey]; !exists {
			return nil, ErrRoleNotFound
		}
		if _, duplicated := seen[roleKey]; duplicated {
			continue
		}
		seen[roleKey] = struct{}{}
		normalized = append(normalized, roleKey)
	}

	admin.Roles = normalized
	admin.UpdatedAt = time.Now()
	return cloneAdmin(admin), nil
}

func (s *MemoryStore) UserByToken(raw string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	token := strings.TrimSpace(strings.TrimPrefix(raw, "Bearer "))
	userID, ok := s.tokens[token]
	if !ok {
		return nil, ErrUnauthorized
	}
	user, ok := s.users[userID]
	if !ok {
		return nil, ErrUserNotFound
	}
	return cloneUser(user), nil
}

func (s *MemoryStore) GetUser(userID uint64) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	user, ok := s.users[userID]
	if !ok {
		return nil, ErrUserNotFound
	}
	return cloneUser(user), nil
}

func (s *MemoryStore) UpdateUser(userID uint64, updater func(*User)) (*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	user, ok := s.users[userID]
	if !ok {
		return nil, ErrUserNotFound
	}
	updater(user)
	user.UpdatedAt = time.Now()
	return cloneUser(user), nil
}

func (s *MemoryStore) ServiceItems() []*ServiceItem {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*ServiceItem, 0, len(s.serviceItems))
	for _, item := range s.serviceItems {
		copyItem := *item
		result = append(result, &copyItem)
	}
	return result
}

func (s *MemoryStore) GetServiceItem(id uint64) (*ServiceItem, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, ok := s.serviceItems[id]
	if !ok {
		return nil, ErrServiceItemNotFound
	}
	copyItem := *item
	return &copyItem, nil
}

func (s *MemoryStore) ApplyProvider(userID uint64, realName, idCardNo string) (*Provider, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if provider, ok := s.providerByUID[userID]; ok {
		provider.RealName = realName
		provider.IDCardNo = idCardNo
		provider.AuditStatus = 1
		provider.AuditRemark = ""
		provider.UpdatedAt = time.Now()
		return cloneProvider(provider), nil
	}

	now := time.Now()
	provider := &Provider{
		ID:          s.nextProviderID,
		UserID:      userID,
		ProviderNo:  fmt.Sprintf("P%08d", s.nextProviderID),
		RealName:    realName,
		IDCardNo:    idCardNo,
		AuditStatus: 1,
		WorkStatus:  1,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	s.nextProviderID++
	s.providers[provider.ID] = provider
	s.providerByUID[userID] = provider
	return cloneProvider(provider), nil
}

func (s *MemoryStore) ProviderByUserID(userID uint64) (*Provider, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	provider, ok := s.providerByUID[userID]
	if !ok {
		return nil, ErrProviderNotFound
	}
	return cloneProvider(provider), nil
}

func (s *MemoryStore) ProviderByID(providerID uint64) (*Provider, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	provider, ok := s.providers[providerID]
	if !ok {
		return nil, ErrProviderNotFound
	}
	return cloneProvider(provider), nil
}

func (s *MemoryStore) Providers() []*Provider {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*Provider, 0, len(s.providers))
	for _, provider := range s.providers {
		result = append(result, cloneProvider(provider))
	}
	return result
}

func (s *MemoryStore) UpdateProviderProfile(userID uint64, displayName, intro string, tags []string) (*Provider, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	provider, ok := s.providerByUID[userID]
	if !ok {
		return nil, ErrProviderNotFound
	}
	provider.DisplayName = displayName
	provider.Intro = intro
	provider.Tags = append([]string(nil), tags...)
	provider.UpdatedAt = time.Now()
	return cloneProvider(provider), nil
}

func (s *MemoryStore) UpdateProviderServiceItems(userID uint64, items []ProviderServiceItem) (*Provider, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	provider, ok := s.providerByUID[userID]
	if !ok {
		return nil, ErrProviderNotFound
	}
	provider.ServiceItems = append([]ProviderServiceItem(nil), items...)
	provider.UpdatedAt = time.Now()
	return cloneProvider(provider), nil
}

func (s *MemoryStore) UpdateProviderWorkStatus(userID uint64, workStatus int) (*Provider, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	provider, ok := s.providerByUID[userID]
	if !ok {
		return nil, ErrProviderNotFound
	}
	provider.WorkStatus = workStatus
	provider.UpdatedAt = time.Now()
	return cloneProvider(provider), nil
}

func (s *MemoryStore) ApproveProvider(providerID uint64, remark string) (*Provider, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	provider, ok := s.providers[providerID]
	if !ok {
		return nil, ErrProviderNotFound
	}
	provider.AuditStatus = 2
	provider.AuditRemark = remark
	provider.UpdatedAt = time.Now()
	return cloneProvider(provider), nil
}

func (s *MemoryStore) RejectProvider(providerID uint64, remark string) (*Provider, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	provider, ok := s.providers[providerID]
	if !ok {
		return nil, ErrProviderNotFound
	}
	provider.AuditStatus = 3
	provider.AuditRemark = remark
	provider.UpdatedAt = time.Now()
	return cloneProvider(provider), nil
}

func (s *MemoryStore) CreateOrder(userID uint64, req CreateOrderInput) (*Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	provider, ok := s.providers[req.ProviderID]
	if !ok {
		return nil, ErrProviderNotFound
	}
	if provider.AuditStatus != 2 {
		return nil, ErrProviderAuditRequired
	}
	item, ok := s.serviceItems[req.ServiceItemID]
	if !ok {
		return nil, ErrServiceItemNotFound
	}
	now := time.Now()
	order := &Order{
		ID:              s.nextOrderID,
		OrderNo:         fmt.Sprintf("L%014d", s.nextOrderID),
		UserID:          userID,
		ProviderID:      req.ProviderID,
		ServiceItemID:   req.ServiceItemID,
		SceneText:       req.SceneText,
		CityCode:        req.CityCode,
		AddressText:     req.AddressText,
		PlannedStartAt:  req.PlannedStartAt,
		PlannedDuration: req.PlannedDuration,
		Status:          OrderStatusPendingPayment,
		PayAmount:       item.MinPrice,
		CreatedAt:       now,
		UpdatedAt:       now,
		Logs: []OrderLog{{
			FromStatus:   0,
			ToStatus:     OrderStatusPendingPayment,
			OperatorRole: "user",
			OperatorID:   userID,
			Remark:       "order created",
			CreatedAt:    now,
		}},
	}
	s.orders[order.ID] = order
	s.nextOrderID++
	return cloneOrder(order), nil
}

func (s *MemoryStore) OrdersByUser(userID uint64) []*Order {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []*Order
	for _, order := range s.orders {
		if order.UserID == userID {
			result = append(result, cloneOrder(order))
		}
	}
	return result
}

func (s *MemoryStore) OrdersByProvider(userID uint64) ([]*Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	provider, ok := s.providerByUID[userID]
	if !ok {
		return nil, ErrProviderNotFound
	}
	var result []*Order
	for _, order := range s.orders {
		if order.ProviderID == provider.ID {
			result = append(result, cloneOrder(order))
		}
	}
	return result, nil
}

func (s *MemoryStore) GetOrder(orderID uint64) (*Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	order, ok := s.orders[orderID]
	if !ok {
		return nil, ErrOrderNotFound
	}
	return cloneOrder(order), nil
}

func (s *MemoryStore) CreatePayment(orderID uint64) (*Payment, *Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	order, ok := s.orders[orderID]
	if !ok {
		return nil, nil, ErrOrderNotFound
	}
	if order.Status != OrderStatusPendingPayment {
		return nil, nil, ErrInvalidOrderStatus
	}
	now := time.Now()
	payment := &Payment{
		ID:        s.nextPaymentID,
		OrderID:   orderID,
		PayAmount: order.PayAmount,
		Status:    20,
		CreatedAt: now,
		PaidAt:    now,
	}
	s.payments[payment.ID] = payment
	s.nextPaymentID++
	s.transitionOrder(order, OrderStatusPendingAccept, "user", order.UserID, "payment completed")
	return payment, cloneOrder(order), nil
}

func (s *MemoryStore) CancelOrder(orderID, userID uint64, reason string) (*Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	order, ok := s.orders[orderID]
	if !ok {
		return nil, ErrOrderNotFound
	}
	if order.UserID != userID {
		return nil, ErrUnauthorized
	}
	switch order.Status {
	case OrderStatusPendingPayment, OrderStatusPendingAccept, OrderStatusAccepted:
	default:
		return nil, ErrInvalidOrderStatus
	}
	order.CancelReason = reason
	s.transitionOrder(order, OrderStatusCanceled, "user", userID, reason)
	return cloneOrder(order), nil
}

func (s *MemoryStore) ProviderAcceptOrder(userID, orderID uint64) (*Order, error) {
	return s.transitionByProvider(userID, orderID, OrderStatusPendingAccept, OrderStatusAccepted, "provider accepted")
}

func (s *MemoryStore) ProviderDepartOrder(userID, orderID uint64) (*Order, error) {
	return s.transitionByProvider(userID, orderID, OrderStatusAccepted, OrderStatusDeparted, "provider departed")
}

func (s *MemoryStore) ProviderArriveOrder(userID, orderID uint64) (*Order, error) {
	return s.transitionByProvider(userID, orderID, OrderStatusDeparted, OrderStatusArrived, "provider arrived")
}

func (s *MemoryStore) StartOrder(userID, orderID uint64, remark string) (*Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	order, ok := s.orders[orderID]
	if !ok {
		return nil, ErrOrderNotFound
	}
	if order.UserID != userID {
		return nil, ErrUnauthorized
	}
	if order.Status != OrderStatusArrived {
		return nil, ErrInvalidOrderStatus
	}
	now := time.Now()
	order.StartedAt = &now
	s.transitionOrder(order, OrderStatusServing, "user", userID, remark)
	return cloneOrder(order), nil
}

func (s *MemoryStore) ProviderFinishOrder(userID, orderID uint64) (*Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	provider, ok := s.providerByUID[userID]
	if !ok {
		return nil, ErrProviderNotFound
	}
	order, ok := s.orders[orderID]
	if !ok {
		return nil, ErrOrderNotFound
	}
	if order.ProviderID != provider.ID {
		return nil, ErrUnauthorized
	}
	if order.Status != OrderStatusServing {
		return nil, ErrInvalidOrderStatus
	}
	now := time.Now()
	order.FinishedAt = &now
	s.transitionOrder(order, OrderStatusPendingFinish, "provider", userID, "provider finished service")
	return cloneOrder(order), nil
}

func (s *MemoryStore) ConfirmFinishOrder(userID, orderID uint64) (*Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	order, ok := s.orders[orderID]
	if !ok {
		return nil, ErrOrderNotFound
	}
	if order.UserID != userID {
		return nil, ErrUnauthorized
	}
	if order.Status != OrderStatusPendingFinish {
		return nil, ErrInvalidOrderStatus
	}
	s.transitionOrder(order, OrderStatusCompleted, "user", userID, "user confirmed finish")
	return cloneOrder(order), nil
}

func (s *MemoryStore) transitionByProvider(userID, orderID uint64, expectedStatus, targetStatus int, remark string) (*Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	provider, ok := s.providerByUID[userID]
	if !ok {
		return nil, ErrProviderNotFound
	}
	order, ok := s.orders[orderID]
	if !ok {
		return nil, ErrOrderNotFound
	}
	if provider.AuditStatus != 2 {
		return nil, ErrProviderAuditRequired
	}
	if order.ProviderID != provider.ID {
		return nil, ErrUnauthorized
	}
	if order.Status != expectedStatus {
		return nil, ErrInvalidOrderStatus
	}
	s.transitionOrder(order, targetStatus, "provider", userID, remark)
	return cloneOrder(order), nil
}

func (s *MemoryStore) transitionOrder(order *Order, toStatus int, role string, operatorID uint64, remark string) {
	now := time.Now()
	order.Logs = append(order.Logs, OrderLog{
		FromStatus:   order.Status,
		ToStatus:     toStatus,
		OperatorRole: role,
		OperatorID:   operatorID,
		Remark:       remark,
		CreatedAt:    now,
	})
	order.Status = toStatus
	order.UpdatedAt = now
}

func cloneUser(user *User) *User {
	if user == nil {
		return nil
	}
	copyUser := *user
	return &copyUser
}

func cloneProvider(provider *Provider) *Provider {
	if provider == nil {
		return nil
	}
	copyProvider := *provider
	copyProvider.Tags = append([]string(nil), provider.Tags...)
	copyProvider.ServiceItems = append([]ProviderServiceItem(nil), provider.ServiceItems...)
	return &copyProvider
}

func cloneOrder(order *Order) *Order {
	if order == nil {
		return nil
	}
	copyOrder := *order
	copyOrder.Logs = append([]OrderLog(nil), order.Logs...)
	return &copyOrder
}
