package model

import (
	"sync"
	"time"
)

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

	now := time.Now()
	defaultAdmin := &AdminUser{ID: s.nextAdminID, Username: "admin", Nickname: "系统管理员", Roles: []string{"super_admin"}, Status: 1, CreatedAt: now, UpdatedAt: now}
	s.nextAdminID++
	s.adminUsers[defaultAdmin.ID] = defaultAdmin
	s.adminByName[defaultAdmin.Username] = defaultAdmin
	s.adminPasswords[defaultAdmin.ID] = "admin123456"

	contentAdmin := &AdminUser{ID: s.nextAdminID, Username: "content_admin", Nickname: "内容管理员", Roles: []string{"content_admin"}, Status: 1, CreatedAt: now, UpdatedAt: now}
	s.nextAdminID++
	s.adminUsers[contentAdmin.ID] = contentAdmin
	s.adminByName[contentAdmin.Username] = contentAdmin
	s.adminPasswords[contentAdmin.ID] = "admin123456"
	return s
}
