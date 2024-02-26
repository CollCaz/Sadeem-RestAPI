CREATE TABLE IF NOT EXISTS user_roles(
  id bigserial PRIMARY KEY,
  user_id bigserial,
  role_id bigserial, 

  CONSTRAINT fk_user
    FOREIGN KEY(user_id)
    REFERENCES users(id),

  CONSTRAINT fk_role
    FOREIGN KEY(role_id)
    REFERENCES roles(id)
);
