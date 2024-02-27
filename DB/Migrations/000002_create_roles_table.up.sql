-- Table storing the possible roles a user can have
CREATE TABLE IF NOT EXISTS roles (
  id bigserial PRIMARY KEY, 
  name VARCHAR(100) UNIQUE NOT NULL,

  CONSTRAINT role_name
    CHECK (name = 'regural' OR name = 'admin')
);

INSERT INTO roles (name) VALUES('regural');
INSERT INTO roles (name) VALUES('admin');
