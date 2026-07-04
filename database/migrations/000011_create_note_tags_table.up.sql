CREATE TABLE note_tags (
    note_id VARCHAR(36) NOT NULL,
    ntag_id VARCHAR(36) NOT NULL,
    PRIMARY KEY (note_id, ntag_id),
    FOREIGN KEY (note_id) REFERENCES notes(id) ON DELETE CASCADE,
    FOREIGN KEY (ntag_id) REFERENCES ntags(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;