CREATE TABLE IF NOT EXISTS user_settings (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id VARCHAR(64) NOT NULL,
  preference_json TEXT NOT NULL,
  notification_json TEXT NOT NULL,
  privacy_json TEXT NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_user_settings_user_id (user_id),
  INDEX idx_user_settings_updated (updated_at)
);

INSERT INTO user_settings(user_id, preference_json, notification_json, privacy_json, created_at, updated_at)
SELECT
  iua.user_id,
  '{"preferSameCityProviders":true,"autoPlaySoundPreview":true,"hideOfflineProviders":false}',
  '{"orderStatusUpdate":true,"complaintResultNotice":true,"marketingActivity":false}',
  '{"profilePublicVisible":true,"personalizedRecommendation":true,"riskControlDataSharing":true}',
  NOW(),
  NOW()
FROM identity_user_accounts iua
LEFT JOIN user_settings us ON us.user_id = iua.user_id
WHERE us.user_id IS NULL;
