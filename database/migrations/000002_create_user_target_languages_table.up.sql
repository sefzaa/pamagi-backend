CREATE TABLE IF NOT EXISTS user_target_languages (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    language_code VARCHAR(10) NOT NULL,  -- Contoh: 'es' untuk Spanyol, 'en' untuk Inggris
    language_name VARCHAR(100) NOT NULL, -- Contoh: 'Spanish', 'English'
    flag_icon VARCHAR(255) NOT NULL,     -- Contoh: Emoji '🇪🇸' atau URL
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;