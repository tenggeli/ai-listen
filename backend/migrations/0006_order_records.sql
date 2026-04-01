CREATE TABLE IF NOT EXISTS order_records (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  order_id VARCHAR(64) NOT NULL,
  user_id VARCHAR(64) NOT NULL,
  provider_id VARCHAR(64) NOT NULL,
  provider_name VARCHAR(128) NOT NULL DEFAULT '',
  service_item_id VARCHAR(64) NOT NULL,
  service_item_title VARCHAR(128) NOT NULL DEFAULT '',
  amount INT NOT NULL,
  currency VARCHAR(16) NOT NULL DEFAULT 'CNY',
  status VARCHAR(32) NOT NULL DEFAULT 'created',
  created_at DATETIME NOT NULL,
  paid_at DATETIME NULL,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_order_records_order_id (order_id),
  INDEX idx_order_records_user_created (user_id, created_at DESC),
  INDEX idx_order_records_status (status)
);
