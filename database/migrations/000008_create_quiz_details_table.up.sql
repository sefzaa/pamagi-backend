CREATE TABLE IF NOT EXISTS quiz_details (
    id VARCHAR(36) PRIMARY KEY,
    quiz_history_id VARCHAR(36) NOT NULL,
    word_id VARCHAR(36) NOT NULL,
    is_correct BOOLEAN DEFAULT NULL,
    FOREIGN KEY (quiz_history_id) REFERENCES quiz_histories(id) ON DELETE CASCADE,
    FOREIGN KEY (word_id) REFERENCES words(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;