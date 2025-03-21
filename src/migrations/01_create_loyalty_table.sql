CREATE TABLE IF NOT EXISTS loyalty (
    id INT auto_increment PRIMARY KEY,
    user_id INT NOT NULL,
    type VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    created BIGINT UNSIGNED,
    active TINYINT DEFAULT 1
);