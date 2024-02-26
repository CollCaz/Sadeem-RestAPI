CREATE TABLE IF NOT EXISTS user_roles(
  id bigserial primary key,
  user_id bigserial references users(id),
  role_id bigserial references roles(id)
);

CREATE OR REPLACE FUNCTION auto_role() RETURNS TRIGGER AS $$
  BEGIN
    if new.roles_id is null then
      new.roles_id = (SELECT id FROM roles WHERE name = 'regural');
    end if;
    return new;
  end

$$ language plpgsql;

  CREATE TRIGGER
      auto_role
    BEFORE INSERT ON 
      user_roles
    FOR EACH ROW EXECUTE PROCEDURE
      auto_role();
