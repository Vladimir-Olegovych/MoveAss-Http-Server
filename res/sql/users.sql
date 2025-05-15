CREATE TABLE IF NOT EXISTS users(
  id INTEGER PRIMARY KEY AUTOINCREMENT, 
  user_token TEXT,
  user_name TEXT,
  user_password TEXT
);