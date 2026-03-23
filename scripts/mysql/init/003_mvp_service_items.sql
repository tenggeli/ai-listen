USE ai_listen;

CREATE TABLE IF NOT EXISTS `service_items` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `provider_id` BIGINT UNSIGNED NOT NULL,
  `category_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `title` VARCHAR(128) NOT NULL,
  `description` VARCHAR(500) DEFAULT '',
  `unit_price` DECIMAL(10,2) NOT NULL DEFAULT 0.00,
  `billing_type` TINYINT NOT NULL DEFAULT 1 COMMENT '1按小时2按次',
  `min_hours` INT NOT NULL DEFAULT 0,
  `max_hours` INT NOT NULL DEFAULT 0,
  `unit_name` VARCHAR(20) DEFAULT '',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '1上架0下架',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` DATETIME DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_provider_id` (`provider_id`),
  KEY `idx_category_id` (`category_id`),
  KEY `idx_status` (`status`),
  KEY `idx_provider_id_status` (`provider_id`, `status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
