CREATE TABLE IF NOT EXISTS admin_users (
  id bigserial PRIMARY KEY, 
  user_id int REFERENCES users(id)
);
