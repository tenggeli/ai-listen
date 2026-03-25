package model

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"ai-listen/backend/migrations"
	_ "github.com/go-sql-driver/mysql"
)

const (
	operatorRoleUser     = 1
	operatorRoleProvider = 2
	operatorRoleAdmin    = 3
)

type MySQLStore struct {
	db           *sql.DB
	adminTokenMu sync.RWMutex
	adminTokens  map[string]uint64
}

func NewMySQLStore(dsn string) (*MySQLStore, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("open mysql: %w", err)
	}

	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping mysql: %w", err)
	}

	store := &MySQLStore{
		db:          db,
		adminTokens: map[string]uint64{},
	}
	if err := store.initSchema(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := store.seedServiceItems(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := store.seedAdminBootstrap(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	return store, nil
}

func (s *MySQLStore) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *MySQLStore) initSchema(ctx context.Context) error {
	for _, stmt := range splitSQLStatements(migrations.InitSchemaSQL) {
		if _, err := s.db.ExecContext(ctx, stmt); err != nil {
			return fmt.Errorf("exec schema: %w", err)
		}
	}
	return nil
}

func (s *MySQLStore) seedServiceItems(ctx context.Context) error {
	var count int
	if err := s.db.QueryRowContext(ctx, "SELECT COUNT(1) FROM service_items").Scan(&count); err != nil {
		return fmt.Errorf("count service_items: %w", err)
	}
	if count > 0 {
		return nil
	}

	items := []struct {
		name     string
		category string
		unit     string
		minPrice int64
		maxPrice int64
		sort     int
	}{
		{name: "陪吃饭", category: "陪伴", unit: "小时", minPrice: 10000, maxPrice: 50000, sort: 10},
		{name: "观影搭子", category: "娱乐", unit: "小时", minPrice: 8000, maxPrice: 40000, sort: 20},
		{name: "心理疏导", category: "情绪支持", unit: "小时", minPrice: 20000, maxPrice: 100000, sort: 30},
	}
	for _, item := range items {
		if _, err := s.db.ExecContext(ctx, `
			INSERT INTO service_items(name, category, unit, min_price, max_price, sort, status)
			VALUES (?, ?, ?, ?, ?, ?, 1)
		`, item.name, item.category, item.unit, item.minPrice, item.maxPrice, item.sort); err != nil {
			return fmt.Errorf("seed service_items: %w", err)
		}
	}
	return nil
}
