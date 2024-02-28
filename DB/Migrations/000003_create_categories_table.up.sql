CREATE TABLE IF NOT EXISTS categories (
  id bigserial primary key,
  name text UNIQUE NOT NULL
);

-- Some test categories
INSERT INTO categories (name) VALUES ('Chairs');
INSERT INTO categories (name) VALUES ('Desks');
INSERT INTO categories (name) VALUES ('Tables');
INSERT INTO categories (name) VALUES ('Sofas');
