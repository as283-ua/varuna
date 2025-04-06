DROP TABLE IF EXISTS users_files CASCADE;
DROP TABLE IF EXISTS files CASCADE;
DROP TABLE IF EXISTS users CASCADE;

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    pubkey TEXT,
    "createdAt" TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    "updatedAt" TIMESTAMPTZ DEFAULT NOW() NOT NULL,
);

CREATE TYPE roles AS ENUM ('Software', 'Hardware', 'HR', 'Finance', 'QA')

CREATE TABLE files (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    user_id BIGINT,
    hash TEXT NOT NULL,
    "createdAt" TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    "updatedAt" TIMESTAMPTZ DEFAULT NOW() NOT NULL,
);

CREATE TABLE users_files (
    user_id BIGINT,
    file_id BIGINT,
    enc_personal_key VARCHAR(376),              -- AES256 encrypted with user's public key, in base64 + the AES key's 16 byte iv in hexadecimal. Encrypting with rsa 256 returns 256 bytes, in base 64 that's 344 + 32 bytes of iv = length 376
    "createdAt" TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    "updatedAt" TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    PRIMARY KEY (user_id, file_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE
);