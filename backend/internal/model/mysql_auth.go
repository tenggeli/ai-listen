package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

func (s *MySQLStore) IssueSMSCode(mobile string) string {
	code := "123456"
	expiresAt := time.Now().Add(10 * time.Minute)
	_, _ = s.db.Exec(`
		INSERT INTO auth_sms_codes(mobile, code, expires_at)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE code = VALUES(code), expires_at = VALUES(expires_at), updated_at = CURRENT_TIMESTAMP
	`, mobile, code, expiresAt)
	_, _ = s.db.Exec(`
		INSERT INTO sms_logs(mobile, scene, template_code, status, provider_msg_id)
		VALUES (?, 'login', 'debug', 20, '')
	`, mobile)
	return code
}

func (s *MySQLStore) LoginBySMS(mobile, code string) (*User, string, string, error) {
	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, "", "", err
	}
	defer tx.Rollback()

	var storedCode string
	var expiresAt time.Time
	err = tx.QueryRowContext(ctx, `
		SELECT code, expires_at
		FROM auth_sms_codes
		WHERE mobile = ?
	`, mobile).Scan(&storedCode, &expiresAt)
	if err != nil || storedCode != code || expiresAt.Before(time.Now()) {
		return nil, "", "", ErrInvalidSMSCode
	}

	user, err := s.findOrCreateUserByMobileTx(ctx, tx, mobile)
	if err != nil {
		return nil, "", "", err
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO user_auths(user_id, auth_type, auth_key)
		VALUES (?, 'sms', ?)
		ON DUPLICATE KEY UPDATE user_id = VALUES(user_id), updated_at = CURRENT_TIMESTAMP
	`, user.ID, mobile); err != nil {
		return nil, "", "", err
	}

	accessToken := fmt.Sprintf("token-%d-%d", user.ID, time.Now().UnixNano())
	refreshToken := fmt.Sprintf("refresh-%d-%d", user.ID, time.Now().UnixNano())
	if err := s.insertTokenTx(ctx, tx, user.ID, accessToken, "access", 24*time.Hour); err != nil {
		return nil, "", "", err
	}
	if err := s.insertTokenTx(ctx, tx, user.ID, refreshToken, "refresh", 7*24*time.Hour); err != nil {
		return nil, "", "", err
	}

	if _, err := tx.ExecContext(ctx, `UPDATE users SET last_login_at = CURRENT_TIMESTAMP WHERE id = ?`, user.ID); err != nil {
		return nil, "", "", err
	}
	if err := tx.Commit(); err != nil {
		return nil, "", "", err
	}
	return user, accessToken, refreshToken, nil
}

func (s *MySQLStore) RefreshToken(refreshToken string) (string, error) {
	token := strings.TrimSpace(refreshToken)
	var userID uint64
	var expiresAt time.Time
	err := s.db.QueryRow(`
		SELECT user_id, expires_at
		FROM auth_tokens
		WHERE token = ? AND token_type = 'refresh' AND status = 1
	`, token).Scan(&userID, &expiresAt)
	if err != nil || expiresAt.Before(time.Now()) {
		return "", ErrUnauthorized
	}

	accessToken := fmt.Sprintf("token-%d-%d", userID, time.Now().UnixNano())
	if _, err := s.db.Exec(`
		INSERT INTO auth_tokens(user_id, token, token_type, expires_at, status)
		VALUES (?, ?, 'access', ?, 1)
	`, userID, accessToken, time.Now().Add(24*time.Hour)); err != nil {
		return "", err
	}
	return accessToken, nil
}
func (s *MySQLStore) insertTokenTx(ctx context.Context, tx *sql.Tx, userID uint64, token, tokenType string, ttl time.Duration) error {
	_, err := tx.ExecContext(ctx, `
		INSERT INTO auth_tokens(user_id, token, token_type, expires_at, status)
		VALUES (?, ?, ?, ?, 1)
	`, userID, token, tokenType, time.Now().Add(ttl))
	return err
}
