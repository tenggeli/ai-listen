package model

import (
	"context"
	"database/sql"
	"fmt"
)

func (s *MySQLStore) ApplyProvider(userID uint64, realName, idCardNo string) (*Provider, error) {
	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var providerID uint64
	err = tx.QueryRowContext(ctx, `SELECT id FROM providers WHERE user_id = ?`, userID).Scan(&providerID)
	switch {
	case err == nil:
		if _, err := tx.ExecContext(ctx, `
			UPDATE providers
			SET real_name = ?, id_card_no = ?, audit_status = 1
			WHERE id = ?
		`, realName, idCardNo, providerID); err != nil {
			return nil, err
		}
	case err == sql.ErrNoRows:
		res, err := tx.ExecContext(ctx, `
			INSERT INTO providers(user_id, provider_no, real_name, id_card_no, audit_status, work_status, status)
			VALUES (?, ?, ?, ?, 1, 1, 1)
		`, userID, fmt.Sprintf("P%08d", userID), realName, idCardNo)
		if err != nil {
			return nil, err
		}
		id, _ := res.LastInsertId()
		providerID = uint64(id)
	default:
		return nil, err
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO audit_records(biz_type, biz_id, audit_status, auditor_id, audit_remark)
		VALUES ('provider', ?, 10, 0, '')
	`, providerID); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return s.ProviderByUserID(userID)
}

func (s *MySQLStore) ProviderByUserID(userID uint64) (*Provider, error) {
	var providerID uint64
	if err := s.db.QueryRow(`SELECT id FROM providers WHERE user_id = ? AND deleted_at IS NULL`, userID).Scan(&providerID); err != nil {
		return nil, ErrProviderNotFound
	}
	return s.ProviderByID(providerID)
}

func (s *MySQLStore) ProviderByID(providerID uint64) (*Provider, error) {
	provider, err := s.loadProvider(providerID)
	if err != nil {
		return nil, err
	}
	return provider, nil
}

func (s *MySQLStore) Providers() []*Provider {
	rows, err := s.db.Query(`SELECT id FROM providers WHERE deleted_at IS NULL ORDER BY id ASC`)
	if err != nil {
		return []*Provider{}
	}
	defer rows.Close()

	var list []*Provider
	for rows.Next() {
		var id uint64
		if err := rows.Scan(&id); err != nil {
			continue
		}
		provider, err := s.loadProvider(id)
		if err == nil {
			list = append(list, provider)
		}
	}
	return list
}

func (s *MySQLStore) UpdateProviderProfile(userID uint64, displayName, intro string, tags []string) (*Provider, error) {
	provider, err := s.ProviderByUserID(userID)
	if err != nil {
		return nil, err
	}
	tagsJSON := toJSON(tags)
	if _, err := s.db.Exec(`
		INSERT INTO provider_profiles(provider_id, display_name, intro, tags)
		VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE display_name = VALUES(display_name), intro = VALUES(intro), tags = VALUES(tags), updated_at = CURRENT_TIMESTAMP
	`, provider.ID, displayName, intro, tagsJSON); err != nil {
		return nil, err
	}
	return s.ProviderByID(provider.ID)
}

func (s *MySQLStore) UpdateProviderServiceItems(userID uint64, items []ProviderServiceItem) (*Provider, error) {
	provider, err := s.ProviderByUserID(userID)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `DELETE FROM provider_service_items WHERE provider_id = ?`, provider.ID); err != nil {
		return nil, err
	}
	for _, item := range items {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO provider_service_items(provider_id, service_item_id, price_amount, price_unit, status)
			VALUES (?, ?, ?, ?, 1)
		`, provider.ID, item.ServiceItemID, item.PriceAmount, item.PriceUnit); err != nil {
			return nil, err
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return s.ProviderByID(provider.ID)
}

func (s *MySQLStore) UpdateProviderWorkStatus(userID uint64, workStatus int) (*Provider, error) {
	provider, err := s.ProviderByUserID(userID)
	if err != nil {
		return nil, err
	}
	if _, err := s.db.Exec(`UPDATE providers SET work_status = ? WHERE id = ?`, workStatus, provider.ID); err != nil {
		return nil, err
	}
	return s.ProviderByID(provider.ID)
}

func (s *MySQLStore) ApproveProvider(providerID uint64, remark string) (*Provider, error) {
	return s.updateProviderAuditStatus(providerID, 2, 20, remark)
}

func (s *MySQLStore) RejectProvider(providerID uint64, remark string) (*Provider, error) {
	return s.updateProviderAuditStatus(providerID, 3, 30, remark)
}
func (s *MySQLStore) updateProviderAuditStatus(providerID uint64, providerStatus, auditRecordStatus int, remark string) (*Provider, error) {
	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(ctx, `UPDATE providers SET audit_status = ? WHERE id = ?`, providerStatus, providerID)
	if err != nil {
		return nil, err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return nil, ErrProviderNotFound
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO audit_records(biz_type, biz_id, audit_status, auditor_id, audit_remark)
		VALUES ('provider', ?, ?, 1, ?)
	`, providerID, auditRecordStatus, remark); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return s.ProviderByID(providerID)
}

func (s *MySQLStore) loadProvider(providerID uint64) (*Provider, error) {
	provider := &Provider{}
	err := s.db.QueryRow(`
		SELECT id, user_id, provider_no, real_name, id_card_no, audit_status, work_status, created_at, updated_at
		FROM providers
		WHERE id = ? AND deleted_at IS NULL
	`, providerID).Scan(
		&provider.ID, &provider.UserID, &provider.ProviderNo, &provider.RealName, &provider.IDCardNo,
		&provider.AuditStatus, &provider.WorkStatus, &provider.CreatedAt, &provider.UpdatedAt,
	)
	if err != nil {
		return nil, ErrProviderNotFound
	}

	var tagsRaw sql.NullString
	_ = s.db.QueryRow(`SELECT display_name, intro, tags FROM provider_profiles WHERE provider_id = ?`, providerID).Scan(&provider.DisplayName, &provider.Intro, &tagsRaw)
	provider.Tags = parseJSONStringArray(tagsRaw.String)
	provider.AuditRemark = s.latestAuditRemark("provider", providerID)
	provider.ServiceItems = s.loadProviderServiceItems(providerID)
	return provider, nil
}

func (s *MySQLStore) latestAuditRemark(bizType string, bizID uint64) string {
	var remark string
	_ = s.db.QueryRow(`
		SELECT audit_remark
		FROM audit_records
		WHERE biz_type = ? AND biz_id = ?
		ORDER BY id DESC
		LIMIT 1
	`, bizType, bizID).Scan(&remark)
	return remark
}

func (s *MySQLStore) loadProviderServiceItems(providerID uint64) []ProviderServiceItem {
	rows, err := s.db.Query(`
		SELECT service_item_id, price_amount, price_unit
		FROM provider_service_items
		WHERE provider_id = ? AND status = 1
		ORDER BY id ASC
	`, providerID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var items []ProviderServiceItem
	for rows.Next() {
		var item ProviderServiceItem
		if err := rows.Scan(&item.ServiceItemID, &item.PriceAmount, &item.PriceUnit); err == nil {
			items = append(items, item)
		}
	}
	return items
}
