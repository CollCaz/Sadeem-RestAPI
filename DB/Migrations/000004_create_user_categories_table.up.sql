CREATE TABLE IF NOT EXISTS user_categories (
  user_id int NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  category_id int NOT NULL REFERENCES categories(id) ON DELETE CASCADE,

  UNIQUE (user_id, category_id)
);
