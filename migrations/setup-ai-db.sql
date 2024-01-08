ALTER DATABASE ai_db SET timezone TO 'Europe/Moscow';
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

ALTER SYSTEM SET wal_level = logical;

-- AI
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