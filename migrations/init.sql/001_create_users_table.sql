CREATE TABLE users
(
    id       SERIAL PRIMARY KEY,
    name     VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    balance  BIGINT       NOT NULL
);