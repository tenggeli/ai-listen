ALTER TABLE order_records
  ADD COLUMN IF NOT EXISTS status_action_reason VARCHAR(255) NOT NULL DEFAULT '' AFTER status,
  ADD COLUMN IF NOT EXISTS status_updated_at DATETIME NULL AFTER status_action_reason;
