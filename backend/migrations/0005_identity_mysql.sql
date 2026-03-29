CREATE TABLE IF NOT EXISTS identity_user_accounts (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id VARCHAR(64) NOT NULL,
  phone VARCHAR(32) NULL,
  wechat_open_id VARCHAR(128) NULL,
  register_source VARCHAR(32) NOT NULL DEFAULT 'unknown',
  status VARCHAR(32) NOT NULL DEFAULT 'active',
  display_name VARCHAR(64) NOT NULL DEFAULT '',
  avatar_url VARCHAR(255) NOT NULL DEFAULT '',
  age_range VARCHAR(32) NOT NULL DEFAULT '',
  city VARCHAR(64) NOT NULL DEFAULT '',
  bio VARCHAR(255) NOT NULL DEFAULT '',
  gender VARCHAR(16) NOT NULL DEFAULT '',
  profile_completed TINYINT(1) NOT NULL DEFAULT 0,
  mbti VARCHAR(8) NOT NULL DEFAULT '',
  interest_tags_json TEXT NOT NULL,
  personality_skipped TINYINT(1) NOT NULL DEFAULT 0,
  personality_completed TINYINT(1) NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_identity_user_id (user_id),
  UNIQUE KEY uk_identity_phone (phone),
  UNIQUE KEY uk_identity_wechat (wechat_open_id),
  INDEX idx_identity_status (status)
);

-- MySQL 5.7 compatibility: avoid JSON default expression and backfill existing rows.
UPDATE identity_user_accounts
SET interest_tags_json = '[]'
WHERE interest_tags_json IS NULL OR interest_tags_json = '';
