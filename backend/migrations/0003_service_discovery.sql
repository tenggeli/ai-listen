ALTER TABLE provider_profiles
  ADD COLUMN IF NOT EXISTS avatar_url VARCHAR(255) NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS verification_text VARCHAR(64) NOT NULL DEFAULT '';

CREATE TABLE IF NOT EXISTS service_categories (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  category_id VARCHAR(64) NOT NULL,
  name VARCHAR(64) NOT NULL,
  icon VARCHAR(64) NOT NULL DEFAULT '',
  sort_order INT NOT NULL DEFAULT 0,
  status VARCHAR(16) NOT NULL DEFAULT 'active',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_service_category_id (category_id)
);

CREATE TABLE IF NOT EXISTS provider_service_items (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  item_id VARCHAR(64) NOT NULL,
  provider_id VARCHAR(64) NOT NULL,
  category_id VARCHAR(64) NOT NULL,
  title VARCHAR(128) NOT NULL,
  description VARCHAR(255) NOT NULL DEFAULT '',
  price_amount INT NOT NULL DEFAULT 0,
  price_unit VARCHAR(32) NOT NULL DEFAULT '',
  support_online TINYINT(1) NOT NULL DEFAULT 0,
  sort_order INT NOT NULL DEFAULT 0,
  service_status VARCHAR(16) NOT NULL DEFAULT 'active',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_provider_service_item (item_id),
  INDEX idx_provider_sort (provider_id, service_status, sort_order),
  INDEX idx_category_status (category_id, service_status)
);

CREATE TABLE IF NOT EXISTS provider_public_tags (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  provider_id VARCHAR(64) NOT NULL,
  tag_name VARCHAR(64) NOT NULL,
  sort_order INT NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_provider_tag (provider_id, sort_order)
);

INSERT INTO providers(id, user_id, provider_status, review_status, rating_avg, order_completed_count)
VALUES
  ('p_pub_001', 'u_pub_001', 'active', 'approved', 4.90, 128),
  ('p_pub_002', 'u_pub_002', 'active', 'approved', 4.80, 76),
  ('p_pub_003', 'u_pub_003', 'active', 'approved', 4.90, 205)
ON DUPLICATE KEY UPDATE
  provider_status = VALUES(provider_status),
  review_status = VALUES(review_status),
  rating_avg = VALUES(rating_avg),
  order_completed_count = VALUES(order_completed_count),
  updated_at = NOW();

INSERT INTO provider_profiles(provider_id, display_name, city_code, bio, avatar_url, verification_text)
VALUES
  ('p_pub_001', '暖心倾听师 · 小林', '310100', '擅长情绪陪伴、轻聊天与晚间低压见面。', 'https://mock.listen.local/avatar/p_pub_001.png', '实名认证'),
  ('p_pub_002', '电影散步搭子 · 念念', '440100', '适合想出门透口气的人，节奏轻松不尴尬。', 'https://mock.listen.local/avatar/p_pub_002.png', '高评分'),
  ('p_pub_003', '睡前放松向导 · 阿乔', '110100', '偏向声音陪伴、睡前放松和安抚式聊天。', 'https://mock.listen.local/avatar/p_pub_003.png', '夜间优先')
ON DUPLICATE KEY UPDATE
  display_name = VALUES(display_name),
  city_code = VALUES(city_code),
  bio = VALUES(bio),
  avatar_url = VALUES(avatar_url),
  verification_text = VALUES(verification_text),
  updated_at = NOW();

INSERT INTO service_categories(category_id, name, icon, sort_order, status)
VALUES
  ('cat_all', '全部', 'sparkles', 1, 'active'),
  ('cat_food', '饭搭子', 'utensils', 2, 'active'),
  ('cat_movie', '电影搭子', 'film', 3, 'active'),
  ('cat_chat', '散步聊天', 'message-circle', 4, 'active'),
  ('cat_relax', '心理疏导', 'heart', 5, 'active')
ON DUPLICATE KEY UPDATE
  name = VALUES(name),
  icon = VALUES(icon),
  sort_order = VALUES(sort_order),
  status = VALUES(status),
  updated_at = NOW();

INSERT INTO provider_service_items(item_id, provider_id, category_id, title, description, price_amount, price_unit, support_online, sort_order, service_status)
VALUES
  ('si_001', 'p_pub_001', 'cat_chat', '线上倾听 30 分钟', '适合情绪安抚、想找人说话、睡前放松。', 99, '30分钟', 1, 1, 'active'),
  ('si_002', 'p_pub_001', 'cat_chat', '同城散步聊天 1 小时', '适合下班后想透口气、想有人一起走一段路。', 158, '小时', 0, 2, 'active'),
  ('si_003', 'p_pub_001', 'cat_movie', '电影陪伴 1 场', '偏轻社交陪伴型，适合慢慢破冰。', 188, '次', 0, 3, 'active'),
  ('si_004', 'p_pub_002', 'cat_movie', '电影散步陪伴 1 场', '看电影后可散步聊天，低压轻社交。', 188, '次', 0, 1, 'active'),
  ('si_005', 'p_pub_002', 'cat_chat', '下班散步聊天 1 小时', '适合今晚想换换心情的人。', 158, '小时', 0, 2, 'active'),
  ('si_006', 'p_pub_003', 'cat_relax', '睡前放松语音 30 分钟', '适合焦虑、失眠、需要被稳稳接住的时刻。', 99, '30分钟', 1, 1, 'active')
ON DUPLICATE KEY UPDATE
  title = VALUES(title),
  description = VALUES(description),
  price_amount = VALUES(price_amount),
  price_unit = VALUES(price_unit),
  support_online = VALUES(support_online),
  sort_order = VALUES(sort_order),
  service_status = VALUES(service_status),
  updated_at = NOW();

DELETE FROM provider_public_tags WHERE provider_id IN ('p_pub_001', 'p_pub_002', 'p_pub_003');

INSERT INTO provider_public_tags(provider_id, tag_name, sort_order)
VALUES
  ('p_pub_001', '深夜聊天', 1),
  ('p_pub_001', '温柔倾听', 2),
  ('p_pub_001', '同城可约', 3),
  ('p_pub_002', '电影', 1),
  ('p_pub_002', '散步', 2),
  ('p_pub_002', 'ENFP', 3),
  ('p_pub_003', '助眠', 1),
  ('p_pub_003', '治愈声音', 2),
  ('p_pub_003', 'INFP', 3);
