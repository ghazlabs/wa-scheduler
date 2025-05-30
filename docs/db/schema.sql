CREATE TABLE IF NOT EXISTS messages (
    id VARCHAR(36) PRIMARY KEY,
    content TEXT NOT NULL,
    recipient_numbers TEXT NOT NULL,
    scheduled_sending_at BIGINT,
    sent_at BIGINT,
    retried_count INT DEFAULT 0,
    status VARCHAR(50),
    reason TEXT DEFAULT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);