CREATE INDEX IF NOT EXISTS idx_users_display_name ON users (display_name);
CREATE INDEX IF NOT EXISTS idx_media_title ON media (media_type, title);
CREATE INDEX IF NOT EXISTS idx_user_media_media ON user_media (media_id);
CREATE INDEX IF NOT EXISTS idx_diary_entry_user ON diary_entry (user_id);
CREATE INDEX IF NOT EXISTS idx_diary_entry_media ON diary_entry (media_id);
