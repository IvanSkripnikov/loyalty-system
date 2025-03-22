CREATE TABLE IF NOT EXISTS loyalty_users (
    id INT auto_increment PRIMARY KEY,
    user_id INT NOT NULL,
    loyalty_id INT NOT NULL,
    active TINYINT DEFAULT 1
);