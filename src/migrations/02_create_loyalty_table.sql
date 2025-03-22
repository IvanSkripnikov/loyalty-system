CREATE TABLE IF NOT EXISTS loyalty (
    id INT auto_increment PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    type_id TINYINT NOT NULL,
    manager_id INT DEFAULT 1,
    created DATETIME(3) DEFAULT NOW(),
    expired DATETIME(3) DEFAULT NOW(),
    data TEXT,
    active TINYINT DEFAULT 1
);