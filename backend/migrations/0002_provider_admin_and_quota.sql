CREATE TABLE IF NOT EXISTS ai_daily_quotas (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id VARCHAR(64) NOT NULL,
  quota_date DATE NOT NULL,
  used_count INT NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_user_date (user_id, quota_date)
);

CREATE TABLE IF NOT EXISTS providers (
  id VARCHAR(64) PRIMARY KEY,
  user_id VARCHAR(64) NOT NULL,
  provider_status VARCHAR(32) NOT NULL DEFAULT 'active',
  review_status VARCHAR(32) NOT NULL DEFAULT 'submitted',
  credit_score INT NOT NULL DEFAULT 100,
  rating_avg DECIMAL(3,2) NOT NULL DEFAULT 5.00,
  order_completed_count INT NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_review_status (review_status, provider_status)
);

CREATE TABLE IF NOT EXISTS provider_profiles (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  provider_id VARCHAR(64) NOT NULL,
  display_name VARCHAR(64) NOT NULL,
  city_code VARCHAR(16) NOT NULL DEFAULT '',
  bio VARCHAR(255) NOT NULL DEFAULT '',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_provider_profile (provider_id)
);

CREATE TABLE IF NOT EXISTS provider_review_records (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  provider_id VARCHAR(64) NOT NULL,
  action VARCHAR(32) NOT NULL,
  reason VARCHAR(255) NOT NULL DEFAULT '',
  operator_id VARCHAR(64) NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_provider_created (provider_id, created_at)
);

INSERT INTO providers(id, user_id, provider_status, review_status)
VALUES
  ('p_001', 'u_001', 'active', 'submitted'),
  ('p_002', 'u_002', 'active', 'under_review'),
  ('p_003', 'u_003', 'active', 'supplement_required')
ON DUPLICATE KEY UPDATE updated_at = NOW();

INSERT INTO provider_profiles(provider_id, display_name, city_code, bio)
VALUES
  ('p_001', '暖心倾听师-小林', '310100', '擅长情绪陪伴'),
  ('p_002', '夜谈伙伴-阿泽', '110100', '擅长夜间聊天'),
  ('p_003', '电影散步搭子-念念', '440100', '擅长线下陪伴')
ON DUPLICATE KEY UPDATE updated_at = NOW();
