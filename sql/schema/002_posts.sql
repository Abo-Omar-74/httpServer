-- +goose Up
CREATE TABLE posts(
  id uuid PRIMARY KEY ,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  body VARCHAR NOT NULL,
  user_id uuid NOT NULL,
  FOREIGN KEY (user_id) REFERENCES 
  users(id) ON DELETE CASCADE 
);
-- +goose Down
DROP TABLE posts;