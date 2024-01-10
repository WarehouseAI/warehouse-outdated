ALTER DATABASE users_db SET timezone TO 'Europe/Moscow';
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

ALTER SYSTEM SET wal_level = logical;

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

