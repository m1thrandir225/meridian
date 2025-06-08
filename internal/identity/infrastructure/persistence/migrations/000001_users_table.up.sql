CREATE TABLE users (
    id UUID PRIMARY KEY NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    version BIGINT NOT NULL DEFAULT 1,
    registartion_time TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_users_user_id ON users (id);

CREATE INDEX idx_users_email ON users (email);

CREATE INDEX idx_users_username ON users (username);
