CREATE TABLE transactions
(
    id          SERIAL PRIMARY KEY,
    sender_id   BIGINT                   NOT NULL,
    receiver_id BIGINT                   NOT NULL,
    amount      BIGINT                   NOT NULL CHECK (amount > 0),
    timestamp   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);