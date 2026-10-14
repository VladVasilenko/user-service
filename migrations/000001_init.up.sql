CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL PRIMARY KEY,
                                     name VARCHAR(100) NOT NULL,
    email VARCHAR(150) UNIQUE NOT NULL
    );

INSERT INTO users (name, email) VALUES
                                    ('Алексей', 'alex@example.com'),
                                    ('Иван', 'ivan@example.com'),
                                    ('Мария', 'maria@example.com')
    ON CONFLICT DO NOTHING;
