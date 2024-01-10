ALTER DATABASE auth_db SET timezone TO 'Europe/Moscow';
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS reset_tokens (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id uuid UNIQUE NOT NULL,
  token VARCHAR(255) NOT NULL,
  expires_at TIMESTAMP DEFAULT now() + INTERVAL '10 minutes' NOT NULL,
  created_at TIMESTAMP DEFAULT now() NOT NULL
);

CREATE TABLE IF NOT EXISTS verification_tokens (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id uuid UNIQUE NOT NULL,
  token VARCHAR(255) NOT NULL,
  expires_at TIMESTAMP DEFAULT now() + INTERVAL '10 minutes' NOT NULL,
  created_at TIMESTAMP DEFAULT now() NOT NULL
);