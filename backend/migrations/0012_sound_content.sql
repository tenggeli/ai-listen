CREATE TABLE IF NOT EXISTS audio_sound_categories (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  category_key VARCHAR(64) NOT NULL,
  label VARCHAR(64) NOT NULL,
  sort_order INT NOT NULL DEFAULT 0,
  status VARCHAR(16) NOT NULL DEFAULT 'active',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_audio_sound_category_key (category_key),
  INDEX idx_audio_sound_category_status (status, sort_order)
);

CREATE TABLE IF NOT EXISTS audio_sound_tracks (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  track_id VARCHAR(64) NOT NULL,
  category_key VARCHAR(64) NOT NULL,
  title VARCHAR(128) NOT NULL,
  play_count_text VARCHAR(64) NOT NULL DEFAULT '',
  duration_text VARCHAR(16) NOT NULL DEFAULT '00:00',
  emoji VARCHAR(16) NOT NULL DEFAULT '',
  author VARCHAR(64) NOT NULL DEFAULT '',
  sort_order INT NOT NULL DEFAULT 0,
  status VARCHAR(16) NOT NULL DEFAULT 'active',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_audio_sound_track_id (track_id),
  INDEX idx_audio_sound_track_filter (status, category_key, sort_order)
);

INSERT INTO audio_sound_categories(category_key, label, sort_order, status)
VALUES
  ('all', '全部', 1, 'active'),
  ('nature', '自然白噪音', 2, 'active'),
  ('sleep', '睡眠引导', 3, 'active'),
  ('meditation', '正念冥想', 4, 'active'),
  ('story', '治愈故事', 5, 'active'),
  ('breath', '呼吸练习', 6, 'active')
ON DUPLICATE KEY UPDATE
  label = VALUES(label),
  sort_order = VALUES(sort_order),
  status = VALUES(status),
  updated_at = NOW();

INSERT INTO audio_sound_tracks(track_id, category_key, title, play_count_text, duration_text, emoji, author, sort_order, status)
VALUES
  ('track-rain', 'nature', '深夜雨声 · 安眠版', '1,248 次播放', '35:00', '🌧', 'listen 治愈声音库', 1, 'active'),
  ('track-wave', 'breath', '海浪呼吸引导 · 4-7-8 法', '856 次播放', '12:00', '🌊', 'listen 治愈声音库', 2, 'active'),
  ('track-forest', 'meditation', '森林清晨 · 入睡前冥想', '2,034 次播放', '20:30', '🍃', 'listen 治愈声音库', 3, 'active'),
  ('track-radio', 'story', '今晚陪你说话 · 深夜电台 Vol.12', '3,671 次播放', '28:15', '🌙', 'listen 深夜电台', 4, 'active'),
  ('track-fire', 'nature', '壁炉暖意 · 冬夜陪伴', '987 次播放', '60:00', '🔥', 'listen 治愈声音库', 5, 'active'),
  ('track-star', 'meditation', '星空冥想 · 放下今天的重量', '4,122 次播放', '18:00', '⭐', 'listen 冥想室', 6, 'active')
ON DUPLICATE KEY UPDATE
  category_key = VALUES(category_key),
  title = VALUES(title),
  play_count_text = VALUES(play_count_text),
  duration_text = VALUES(duration_text),
  emoji = VALUES(emoji),
  author = VALUES(author),
  sort_order = VALUES(sort_order),
  status = VALUES(status),
  updated_at = NOW();
