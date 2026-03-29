CREATE TABLE IF NOT EXISTS provider_presence (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  provider_id VARCHAR(64) NOT NULL,
  last_seen_at DATETIME NOT NULL,
  source VARCHAR(32) NOT NULL DEFAULT 'mock_gateway',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_provider_presence (provider_id),
  INDEX idx_last_seen (last_seen_at)
);

INSERT INTO provider_presence(provider_id, last_seen_at, source)
VALUES
  ('p_pub_001', NOW(), 'mock_gateway'),
  ('p_pub_002', DATE_SUB(NOW(), INTERVAL 2 MINUTE), 'mock_gateway'),
  ('p_pub_003', DATE_SUB(NOW(), INTERVAL 8 MINUTE), 'mock_gateway')
ON DUPLICATE KEY UPDATE
  last_seen_at = VALUES(last_seen_at),
  source = VALUES(source),
  updated_at = NOW();
