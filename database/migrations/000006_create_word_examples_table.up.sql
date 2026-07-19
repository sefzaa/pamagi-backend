
CREATE TABLE IF NOT EXISTS word_examples (
    id VARCHAR(36) PRIMARY KEY,
    word_id VARCHAR(36) NOT NULL,
    target_sentence TEXT NOT NULL, -- Kalimat dalam bahasa asing
    native_sentence TEXT NOT NULL, -- Kalimat dalam bahasa ibu
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    FOREIGN KEY (word_id) REFERENCES words(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;