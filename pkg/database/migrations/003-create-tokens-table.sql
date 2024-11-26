-- Create tokens table
CREATE TABLE IF NOT EXISTS tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token TEXT NOT NULL UNIQUE
);

-- Create index on token for faster lookups
CREATE INDEX IF NOT EXISTS tokens_token_idx ON tokens(token);

CREATE INDEX IF NOT EXISTS tokens_user_id_idx ON tokens(user_id);