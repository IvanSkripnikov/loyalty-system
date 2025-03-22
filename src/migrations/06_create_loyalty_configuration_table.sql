CREATE TABLE IF NOT EXISTS loyalty_configuration (
    id INT auto_increment PRIMARY KEY,
    property VARCHAR(255) NOT NULL,
    value VARCHAR(255) NOT NULL,
    active TINYINT DEFAULT 1
);