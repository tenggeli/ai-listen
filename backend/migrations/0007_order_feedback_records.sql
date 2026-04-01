CREATE TABLE IF NOT EXISTS order_feedback_records (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  feedback_id VARCHAR(64) NOT NULL,
  order_id VARCHAR(64) NOT NULL,
  user_id VARCHAR(64) NOT NULL,
  rating_score INT NOT NULL DEFAULT 0,
  review_tags_json TEXT NOT NULL,
  review_content TEXT NOT NULL,
  has_complaint TINYINT(1) NOT NULL DEFAULT 0,
  complaint_reason VARCHAR(128) NOT NULL DEFAULT '',
  complaint_content TEXT NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_feedback_id (feedback_id),
  UNIQUE KEY uk_feedback_order_id (order_id),
  INDEX idx_feedback_user_created (user_id, created_at DESC),
  INDEX idx_feedback_has_complaint (has_complaint)
);

UPDATE order_feedback_records
SET review_tags_json = '[]'
WHERE review_tags_json IS NULL OR review_tags_json = '';
