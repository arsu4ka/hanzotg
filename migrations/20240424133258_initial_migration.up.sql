BEGIN;

CREATE TABLE users
(
    id       VARCHAR(20) PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255)
);

CREATE TYPE verification_message_status AS ENUM (
    'pending',
    'accepted',
    'declined'
);

CREATE TABLE verification_messages
(
    message_id    VARCHAR(20) PRIMARY KEY,
    chat_id       VARCHAR(20) NOT NULL,
    about_user_id VARCHAR(20) REFERENCES users (id),
    status        verification_message_status NOT NULL DEFAULT 'pending'
);

COMMIT;