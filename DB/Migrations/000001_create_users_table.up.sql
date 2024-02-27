CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY, 
    name text UNIQUE NOT NULL,
    email citext UNIQUE NOT NULL,
    profile_picture_path text UNIQUE NOT NULL,
    hashed_password bytea NOT NULL,
    created timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1
);

CREATE OR REPLACE FUNCTION create_role() RETURNS TRIGGER AS $$
  BEGIN
    INSERT INTO user_roles (id) VALUES (new.id);
    return new;
  END

$$ language plpgsql;

 CREATE TRIGGER
     create_role
   AFTER INSERT ON 
     users
   FOR EACH ROW EXECUTE PROCEDURE
     create_role();
