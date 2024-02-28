CREATE TABLE IF NOT EXISTS user_categories (
  user_id int NOT NULL REFERENCES users(id),
  categorie_id int NOT NULL REFERENCES categories(id)
);
