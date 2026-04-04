ALTER TABLE order_records
  ADD COLUMN  status_action_reason VARCHAR(255) NOT NULL DEFAULT '' AFTER status,
  ADD COLUMN status_updated_at DATETIME NULL AFTER status_action_reason;
