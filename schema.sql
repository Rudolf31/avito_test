CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS team (
    id VARCHAR(36) PRIMARY KEY DEFAULT (gen_random_uuid()::VARCHAR(36)),
    team_name VARCHAR NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS "user" (
    id VARCHAR(36) PRIMARY KEY,
    username VARCHAR NOT NULL UNIQUE,
    team_id VARCHAR(36) REFERENCES team(id) NOT NULL,
    is_active BOOLEAN NOT NULL
);

CREATE TABLE IF NOT EXISTS pull_request (
    id VARCHAR(36) PRIMARY KEY,
    pull_request_name VARCHAR UNIQUE NOT NULL,
    author_id VARCHAR(36) NOT NULL,
    status SMALLINT CHECK(status IN (0, 1)) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    merged_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS review (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) REFERENCES "user"(id),
    pull_request_id VARCHAR(36) REFERENCES pull_request(id),
    reviewed BOOLEAN
);

ALTER TABLE pull_request
  ADD CONSTRAINT fk_pull_request_author FOREIGN KEY (author_id) REFERENCES "user"(id);


CREATE INDEX IF NOT EXISTS idx_user_team_id ON "user"(team_id);
CREATE INDEX IF NOT EXISTS idx_pull_request_author_id ON pull_request(author_id);
CREATE INDEX IF NOT EXISTS idx_pull_request_status ON pull_request(status);
CREATE INDEX IF NOT EXISTS idx_team_name ON team(team_name);
CREATE INDEX IF NOT EXISTS idx_user_username ON "user"(username);
CREATE INDEX IF NOT EXISTS idx_review_user_id ON review(user_id);
CREATE INDEX IF NOT EXISTS idx_review_pull_request_id ON review(pull_request_id);