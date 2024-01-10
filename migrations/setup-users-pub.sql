ALTER DATABASE users_db SET timezone TO 'Europe/Moscow';
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE PUBLICATION users_pub FOR All TABLES;