CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    login TEXT UNIQUE,
    password TEXT,
    addition INT DEFAULT 5,
    subtraction INT DEFAULT 5,
    multiplication INT DEFAULT 5,
    division INT DEFAULT 5
);

CREATE TABLE expressions (
    id SERIAL PRIMARY KEY,
    expression TEXT,
    result TEXT DEFAULT '',
    status TEXT,
    created_at TIMESTAMP,
    calculated_at TIMESTAMP,
    calculated_by INT,
    owner_id INT
);

CREATE TABLE cresources (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    address TEXT UNIQUE NOT NULL,
    expression TEXT DEFAULT '',
    occupied BOOL DEFAULT false,
    orchestrator_alive BOOL DEFAULT false
);