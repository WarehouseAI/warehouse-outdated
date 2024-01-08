ALTER DATABASE stat_db SET timezone TO 'Europe/Moscow';
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

ALTER SYSTEM SET wal_level = logical;
ALTER SYSTEM SET max_logical_replication_workers = 5;

/*USERS DB*/

CREATE TABLE IF NOT EXISTS users (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  firstname VARCHAR(255) NOT NULL,
  lastname VARCHAR(255) NOT NULL,
  username VARCHAR(255) NOT NULL UNIQUE,
  picture VARCHAR(255),
  password VARCHAR(72) NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE,
  verified BOOLEAN NOT NULL DEFAULT FALSE,
  via_google BOOLEAN NOT NULL DEFAULT FALSE,
  is_developer BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT now() NOT NULL,
  updated_at TIMESTAMP DEFAULT now() NOT NULL
);

CREATE TABLE IF NOT EXISTS user_favorites (
  id INTEGER PRIMARY KEY,
  ai_id uuid NOT NULL,
  user_id uuid NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS user_owned (
  id INTEGER PRIMARY KEY,
  ai_id uuid NOT NULL,
  user_id uuid NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE OR REPLACE FUNCTION update_updated_at_users()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE
    ON
        users
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_users();

/*AI DB*/

CREATE TABLE IF NOT EXISTS ai_products (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  owner uuid NOT NULL,
  name TEXT NOT NULL,
  description TEXT NOT NULL,
  auth_header_content VARCHAR(255) NOT NULL,
  auth_header_name VARCHAR(40) NOT NULL,
  used INTEGER NOT NULL DEFAULT 0,
  background_url VARCHAR(255) NOT NULL,
  updated_at TIMESTAMP DEFAULT now() NOT NULL,
  created_at TIMESTAMP DEFAULT now() NOT NULL
);

CREATE OR REPLACE FUNCTION update_updated_at_ai_product()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_ai_updated_at
    BEFORE UPDATE
    ON
        ai_products
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_ai_product();

-- COMMANDS
CREATE TABLE IF NOT EXISTS ai_commands (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  ai_id uuid REFERENCES ai_products(id) ON DELETE CASCADE,
  name VARCHAR(255) NOT NULL,
  payload json NOT NULL,
  payload_type VARCHAR(25) NOT NULL,
  request_type VARCHAR(10) NOT NULL,
  input_type VARCHAR(10) NOT NULL,
  output_type VARCHAR(10) NOT NULL,
  url VARCHAR(255) NOT NULL,
  created_at TIMESTAMP DEFAULT now() NOT NULL,
  updated_at TIMESTAMP DEFAULT now() NOT NULL
);

CREATE OR REPLACE FUNCTION update_updated_at_ai_command()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_command_updated_at
    BEFORE UPDATE
    ON
        ai_commands
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_ai_command();

-- RATING
CREATE TABLE IF NOT EXISTS ai_rates (
  id INTEGER PRIMARY KEY,
  by_user_id uuid NOT NULL,
  ai_id uuid NOT NULL,
  rate INTEGER CHECK (rate >= 0 AND rate <= 5),
  created_at TIMESTAMP DEFAULT now() NOT NULL,
  updated_at TIMESTAMP DEFAULT now() NOT NULL
);

CREATE OR REPLACE FUNCTION update_updated_at_ai_rate()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_rate_updated_at
    BEFORE UPDATE
    ON
        ai_rates
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_ai_rate();