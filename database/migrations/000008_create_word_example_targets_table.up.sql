CREATE TABLE IF NOT EXISTS word_example_targets (
    id VARCHAR(36) PRIMARY KEY,
    word_example_id VARCHAR(36) NOT NULL,
    language_code VARCHAR(10) NOT NULL,        -- Penanda bahasa (contoh: 'ru', 'en')
    target_sentence TEXT NOT NULL,             -- Contoh: 'У меня есть часы.', 'I have a clock.'
    FOREIGN KEY (word_example_id) REFERENCES word_examples(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;