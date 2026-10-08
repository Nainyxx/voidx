CREATE TABLE IF NOT EXISTS user_profiles (
    user_id          UUID PRIMARY KEY,
    username         VARCHAR(32)  NOT NULL UNIQUE,
    name             VARCHAR(64)  NOT NULL,
    surname          VARCHAR(64)  NOT NULL DEFAULT '',
    phone            VARCHAR(16)  NOT NULL UNIQUE,
    description      VARCHAR(255) NOT NULL DEFAULT '',
    avatar_image_url TEXT         NOT NULL DEFAULT '',
    created_at       TIMESTAMPTZ  NOT NULL,
    updated_at       TIMESTAMPTZ  NOT NULL
);
