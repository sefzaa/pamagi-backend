CREATE TABLE IF NOT EXISTS quiz_histories (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    total_questions INT NOT NULL,       -- Total Y (contoh: 20)
    correct_answers INT DEFAULT 0,      -- Total X (contoh: 15)
    incorrect_answers INT DEFAULT 0,
    score FLOAT DEFAULT 0,              -- Persentase (contoh: 75.0)
    status ENUM('IN_PROGRESS', 'COMPLETED') DEFAULT 'IN_PROGRESS',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;