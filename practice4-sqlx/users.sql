CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    balance NUMERIC(10, 2) DEFAULT 0
);

INSERT INTO users (name, email, balance) VALUES
('Alice', 'alice@example.com', 100.00),
('Bob', 'bob@example.com', 200.00),
('Charlie', 'charlie@example.com', 50.00);
