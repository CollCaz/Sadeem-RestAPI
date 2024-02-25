CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY, 
    name VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    hashed_password CHAR(60) NOT NULL,
    created timestamp(0) with time zone NOT NULL DEFAULT NOW()
);
