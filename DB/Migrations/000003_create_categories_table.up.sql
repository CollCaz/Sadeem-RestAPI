CREATE TABLE IF NOT EXISTS categories (
  id bigserial primary key,
  name text UNIQUE NOT NULL,
  activated bool DEFAULT FALSE
);

INSERT INTO categories (name) VALUES ('Chairs');
INSERT INTO categories (name) VALUES ('Desks');
INSERT INTO categories (name) VALUES ('Tables');
INSERT INTO categories (name) VALUES ('Sofas');

UPDATE categories SET activated = TRUE WHERE name = ('Sofas');
