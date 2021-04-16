CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id uuid DEFAULT uuid_generate_v4(),
    username VARCHAR UNIQUE NOT NULL,
    password VARCHAR NOT NULL
);