package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"

	domain "listen/backend/internal/domain/identity"
)

type IdentityRepository struct {
	db *sql.DB
}

func NewIdentityRepository(db *sql.DB) IdentityRepository {
	return IdentityRepository{db: db}
}

func (r IdentityRepository) GetByID(ctx context.Context, userID string) (domain.UserAccount, bool, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return domain.UserAccount{}, false, nil
	}

	const queryByUserID = `
SELECT user_id, phone, wechat_open_id, register_source, status, display_name, avatar_url,
       age_range, city, bio, gender, profile_completed, mbti, interest_tags_json,
       personality_skipped, personality_completed
FROM identity_user_accounts
WHERE user_id = ? LIMIT 1`

	account, found, err := r.queryOne(ctx, queryByUserID, userID)
	if err != nil {
		return domain.UserAccount{}, false, err
	}
	if found {
		return account, true, nil
	}

	// Backward-compatible fallback for environments that pass username (display_name)
	// such as "user_1722" instead of canonical user_id.
	const queryByDisplayName = `
SELECT user_id, phone, wechat_open_id, register_source, status, display_name, avatar_url,
       age_range, city, bio, gender, profile_completed, mbti, interest_tags_json,
       personality_skipped, personality_completed
FROM identity_user_accounts
WHERE display_name = ? LIMIT 1`

	account, found, err = r.queryOne(ctx, queryByDisplayName, userID)
	if err != nil {
		return domain.UserAccount{}, false, err
	}
	if !found {
		return domain.UserAccount{}, false, nil
	}
	return account, true, nil
}

func (r IdentityRepository) GetByPhone(ctx context.Context, phone string) (domain.UserAccount, bool, error) {
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return domain.UserAccount{}, false, nil
	}
	const query = `
SELECT user_id, phone, wechat_open_id, register_source, status, display_name, avatar_url,
       age_range, city, bio, gender, profile_completed, mbti, interest_tags_json,
       personality_skipped, personality_completed
FROM identity_user_accounts
WHERE phone = ? LIMIT 1`

	account, found, err := r.queryOne(ctx, query, phone)
	if err != nil {
		return domain.UserAccount{}, false, err
	}
	if !found {
		return domain.UserAccount{}, false, nil
	}
	return account, true, nil
}

func (r IdentityRepository) GetByWechatOpenID(ctx context.Context, openID string) (domain.UserAccount, bool, error) {
	openID = strings.TrimSpace(openID)
	if openID == "" {
		return domain.UserAccount{}, false, nil
	}
	const query = `
SELECT user_id, phone, wechat_open_id, register_source, status, display_name, avatar_url,
       age_range, city, bio, gender, profile_completed, mbti, interest_tags_json,
       personality_skipped, personality_completed
FROM identity_user_accounts
WHERE wechat_open_id = ? LIMIT 1`

	account, found, err := r.queryOne(ctx, query, openID)
	if err != nil {
		return domain.UserAccount{}, false, err
	}
	if !found {
		return domain.UserAccount{}, false, nil
	}
	return account, true, nil
}

func (r IdentityRepository) Save(ctx context.Context, account domain.UserAccount) error {
	tagsJSON, err := json.Marshal(account.InterestTags)
	if err != nil {
		return err
	}

	const upsert = `
INSERT INTO identity_user_accounts(
  user_id, phone, wechat_open_id, register_source, status, display_name, avatar_url,
  age_range, city, bio, gender, profile_completed, mbti, interest_tags_json,
  personality_skipped, personality_completed, created_at, updated_at
) VALUES(
  ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW()
)
ON DUPLICATE KEY UPDATE
  phone = VALUES(phone),
  wechat_open_id = VALUES(wechat_open_id),
  register_source = VALUES(register_source),
  status = VALUES(status),
  display_name = VALUES(display_name),
  avatar_url = VALUES(avatar_url),
  age_range = VALUES(age_range),
  city = VALUES(city),
  bio = VALUES(bio),
  gender = VALUES(gender),
  profile_completed = VALUES(profile_completed),
  mbti = VALUES(mbti),
  interest_tags_json = VALUES(interest_tags_json),
  personality_skipped = VALUES(personality_skipped),
  personality_completed = VALUES(personality_completed),
  updated_at = NOW()`

	_, err = r.db.ExecContext(
		ctx,
		upsert,
		account.UserID,
		toNullString(account.Phone),
		toNullString(account.WechatOpenID),
		account.RegisterSource,
		account.Status,
		account.DisplayName,
		account.AvatarURL,
		account.AgeRange,
		account.City,
		account.Bio,
		account.Gender,
		boolToInt(account.ProfileCompleted),
		account.MBTI,
		string(tagsJSON),
		boolToInt(account.PersonalitySkipped),
		boolToInt(account.PersonalityCompleted),
	)
	return err
}

func (r IdentityRepository) queryOne(ctx context.Context, query string, arg string) (domain.UserAccount, bool, error) {
	var account domain.UserAccount
	var phone sql.NullString
	var openID sql.NullString
	var tagsJSON string
	var profileCompleted int
	var personalitySkipped int
	var personalityCompleted int

	err := r.db.QueryRowContext(ctx, query, arg).Scan(
		&account.UserID,
		&phone,
		&openID,
		&account.RegisterSource,
		&account.Status,
		&account.DisplayName,
		&account.AvatarURL,
		&account.AgeRange,
		&account.City,
		&account.Bio,
		&account.Gender,
		&profileCompleted,
		&account.MBTI,
		&tagsJSON,
		&personalitySkipped,
		&personalityCompleted,
	)
	if err == sql.ErrNoRows {
		return domain.UserAccount{}, false, nil
	}
	if err != nil {
		return domain.UserAccount{}, false, err
	}

	account.Phone = nullToString(phone)
	account.WechatOpenID = nullToString(openID)
	account.ProfileCompleted = profileCompleted == 1
	account.PersonalitySkipped = personalitySkipped == 1
	account.PersonalityCompleted = personalityCompleted == 1
	account.InterestTags = parseTagsJSON(tagsJSON)
	return account, true, nil
}

func toNullString(value string) any {
	v := strings.TrimSpace(value)
	if v == "" {
		return nil
	}
	return v
}

func nullToString(value sql.NullString) string {
	if !value.Valid {
		return ""
	}
	return value.String
}

func boolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func parseTagsJSON(value string) []string {
	value = strings.TrimSpace(value)
	if value == "" {
		return []string{}
	}
	tags := make([]string, 0)
	if err := json.Unmarshal([]byte(value), &tags); err != nil {
		return []string{}
	}
	return tags
}

var _ domain.Repository = IdentityRepository{}
