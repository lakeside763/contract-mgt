ALTER TABLE users DROP CONSTRAINT uni_users_username;

DROP TABLE users;

DROP EXTENSION IF EXISTS "uuid-ossp";