-- +goose Up
CREATE TABLE refresh_tokens(
  token VARCHAR PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  user_id uuid NOT NULL, 
  expires_at TIMESTAMP NOT NULL,
  revoked_at TIMESTAMP DEFAULT NULL,
  FOREIGN KEY (user_id) REFERENCES 
  users(id) ON DELETE CASCADE
);
-- +goose Down
DROP TABLE refresh_tokens;