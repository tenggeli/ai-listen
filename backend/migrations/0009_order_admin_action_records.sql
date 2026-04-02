CREATE TABLE IF NOT EXISTS order_admin_action_records (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  action_id VARCHAR(64) NOT NULL,
  order_id VARCHAR(64) NOT NULL,
  scope VARCHAR(32) NOT NULL DEFAULT 'order',
  action VARCHAR(64) NOT NULL,
  operator_id VARCHAR(64) NOT NULL,
  reason VARCHAR(255) NOT NULL DEFAULT '',
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_order_admin_action_id (action_id),
  INDEX idx_order_admin_action_order_created (order_id, created_at DESC),
  INDEX idx_order_admin_action_scope (scope)
);
