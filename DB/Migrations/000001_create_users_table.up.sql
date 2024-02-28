CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY, 
    name text UNIQUE NOT NULL,
    email citext UNIQUE NOT NULL,
    profile_picture_path text NOT NULL,
    hashed_password bytea NOT NULL,
    created timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1
);
