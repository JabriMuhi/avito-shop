CREATE TABLE purchases
(
    id         SERIAL PRIMARY KEY,
    user_id    BIGINT                   NOT NULL,
    merch_name VARCHAR(255)             NOT NULL,
    timestamp  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);