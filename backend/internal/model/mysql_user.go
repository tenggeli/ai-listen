package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

func (s *MySQLStore) UserByToken(raw string) (*User, error) {
	token := strings.TrimSpace(strings.TrimPrefix(raw, "Bearer "))
	user := &User{}
	var birthday sql.NullTime
	err := s.db.QueryRow(`
		SELECT u.id, u.mobile, u.nickname, u.avatar, u.gender, u.birthday, u.city_code, u.created_at, u.updated_at
		FROM auth_tokens t
		JOIN users u ON u.id = t.user_id
		WHERE t.token = ? AND t.token_type = 'access' AND t.status = 1 AND t.expires_at > CURRENT_TIMESTAMP
	`, token).Scan(
		&user.ID, &user.Mobile, &user.Nickname, &user.Avatar, &user.Gender, &birthday, &user.CityCode, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, ErrUnauthorized
	}
	user.Birthday = formatNullDate(birthday)
	return user, nil
}

func (s *MySQLStore) GetUser(userID uint64) (*User, error) {
	return s.getUser(userID)
}

func (s *MySQLStore) UpdateUser(userID uint64, updater func(*User)) (*User, error) {
	user, err := s.getUser(userID)
	if err != nil {
		return nil, err
	}
	updater(user)

	var birthday any
	if user.Birthday != "" {
		birthday = user.Birthday
	}

	if _, err := s.db.Exec(`
		UPDATE users
		SET nickname = ?, avatar = ?, gender = ?, birthday = ?, city_code = ?
		WHERE id = ?
	`, user.Nickname, user.Avatar, user.Gender, birthday, user.CityCode, userID); err != nil {
		return nil, err
	}
	return s.getUser(userID)
}
func (s *MySQLStore) getUser(userID uint64) (*User, error) {
	user := &User{}
	var birthday sql.NullTime
	err := s.db.QueryRow(`
		SELECT id, mobile, nickname, avatar, gender, birthday, city_code, created_at, updated_at
		FROM users
		WHERE id = ? AND deleted_at IS NULL
	`, userID).Scan(&user.ID, &user.Mobile, &user.Nickname, &user.Avatar, &user.Gender, &birthday, &user.CityCode, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, ErrUserNotFound
	}
	user.Birthday = formatNullDate(birthday)
	return user, nil
}

func (s *MySQLStore) findOrCreateUserByMobileTx(ctx context.Context, tx *sql.Tx, mobile string) (*User, error) {
	user := &User{}
	var birthday sql.NullTime
	err := tx.QueryRowContext(ctx, `
		SELECT id, mobile, nickname, avatar, gender, birthday, city_code, created_at, updated_at
		FROM users
		WHERE mobile = ? AND deleted_at IS NULL
	`, mobile).Scan(&user.ID, &user.Mobile, &user.Nickname, &user.Avatar, &user.Gender, &birthday, &user.CityCode, &user.CreatedAt, &user.UpdatedAt)
	switch {
	case err == nil:
		user.Birthday = formatNullDate(birthday)
		return user, nil
	case err != sql.ErrNoRows:
		return nil, err
	}

	res, err := tx.ExecContext(ctx, `
		INSERT INTO users(mobile, nickname, city_code, status)
		VALUES (?, ?, '310100', 1)
	`, mobile, fmt.Sprintf("listen用户%s", mobile[len(mobile)-4:]))
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	user.ID = uint64(id)
	user.Mobile = mobile
	user.Nickname = fmt.Sprintf("listen用户%s", mobile[len(mobile)-4:])
	user.CityCode = "310100"
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	return user, nil
}
