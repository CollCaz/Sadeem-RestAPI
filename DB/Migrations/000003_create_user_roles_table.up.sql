-- Table storing the relationship between each user and their role
CREATE TABLE IF NOT EXISTS user_roles(
  id bigserial primary key,
  user_id bigserial references users(id),
  role_id bigserial references roles(id)
);

-- Sets the defaul role_id to be the id of the regural role
CREATE OR REPLACE FUNCTION defaule_role() RETURNS TRIGGER AS $$
  BEGIN
    if new.role_id is null then
      new.role_id = (SELECT id FROM roles WHERE name = 'regural');
    end if;
    return new;
  END

$$ language plpgsql;

  -- Runs the default_role() function before any insert
  CREATE TRIGGER
      auto_role
    BEFORE INSERT ON 
      user_roles
    FOR EACH ROW EXECUTE PROCEDURE
      auto_role();
