
CREATE TABLE IF NOT EXISTS words (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    target_language_code VARCHAR(10) NOT NULL, -- Penanda bahasa (harus match dengan user_target_languages)
    target_word VARCHAR(255) NOT NULL,         -- Kata yang dipelajari
    native_word VARCHAR(255) NOT NULL,         -- Terjemahan ke bahasa ibu
    part_of_speech ENUM('NOUN', 'VERB', 'ADJECTIVE', 'ADVERB', 'PRONOUN', 'PREPOSITION', 'CONJUNCTION', 'INTERJECTION', 'IDIOM') NOT NULL,
    is_favorite BOOLEAN DEFAULT FALSE,
    is_bookmarked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;