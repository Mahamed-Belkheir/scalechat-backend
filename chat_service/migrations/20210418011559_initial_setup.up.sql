CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE chats (
    id uuid DEFAULT uuid_generate_v4(),
    name VARCHAR UNIQUE NOT NULL,
    user_id VARCHAR NOT NULL
);