CREATE TABLE IF NOT EXISTS roles (
  id bigserial PRIMARY KEY, 
  name VARCHAR(100) UNIQUE NOT NULL,

  CONSTRAINT role_name
    CHECK (name = 'regural' OR name = 'admin')
);