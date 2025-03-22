CREATE TABLE IF NOT EXISTS loyalty_type (
    id INT auto_increment PRIMARY KEY,
    title INT NOT NULL,
    description TEXT NOT NULL,
    created DATETIME(3) DEFAULT NOW(),
    updated DATETIME(3) DEFAULT NOW(),
    active TINYINT DEFAULT 1
);