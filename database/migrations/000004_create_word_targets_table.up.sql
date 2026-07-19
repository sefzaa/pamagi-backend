CREATE TABLE IF NOT EXISTS word_targets (
    id VARCHAR(36) PRIMARY KEY,
    word_id VARCHAR(36) NOT NULL,
    language_code VARCHAR(10) NOT NULL,        -- Penanda bahasa (contoh: 'ru', 'en')
    target_word VARCHAR(255) NOT NULL,         -- Kata yang dipelajari (contoh: 'часы', 'clock')
    FOREIGN KEY (word_id) REFERENCES words(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;