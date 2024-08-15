-- +migrate Up
CREATE TABLE IF NOT EXISTS users(
    id UUID NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    public_key VARCHAR(255) NOT NULL,
    access_token VARCHAR(255) NOT NULL,
    backup_token VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);
CREATE TABLE IF NOT EXISTS clients(
    id UUID NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id UUID CONSTRAINT fk_vaults_user_id REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,
    full_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);
CREATE TABLE IF NOT EXISTS vaults(
    id UUID NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    pivot_id UUID CONSTRAINT fk_vaults_pivot_id REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    credential VARCHAR(1023) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS user_vault_pivots(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    user_id UUID CONSTRAINT fk_user_vault_pivots_user_id REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,
    vaults_id UUID CONSTRAINT fk_user_vault_pivots_vaults_id REFERENCES vaults(id) ON UPDATE CASCADE ON DELETE CASCADE
);

-- +migrate Down
DROP TABLE IF EXISTS user_vault_pivots;
DROP TABLE IF EXISTS vaults;
DROP TABLE IF EXISTS clients;
DROP TABLE IF EXISTS users;
