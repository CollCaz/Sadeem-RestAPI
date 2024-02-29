CREATE TABLE IF NOT EXISTS user_categories (
  user_id int NOT NULL REFERENCES users(id),
  category_id int NOT NULL REFERENCES categories(id),

  UNIQUE (user_id, category_id)
);
