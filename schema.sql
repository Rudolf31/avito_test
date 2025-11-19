CREATE TYPE IF NOT EXISTS pr_status AS ENUM ('open', 'merged');

CREATE TABLE IF NOT EXISTS pull_request (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pull_request_name VARCHAR UNIQUE NOT NULL,
    author_id UUID REFERENCES user(id) NOT NULL,
    status pr_status NOT NULL,
    created_at TIMESTAMP NOT NULL,
    merged_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS team (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    team_name VARCHAR NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS "user" (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR NOT NULL UNIQUE,
    team_id UUID REFERENCES team(id),
    is_active BOOLEAN NOT NULL,
    assigned_pull_request_id UUID REFERENCES pull_request(id)
);

CREATE INDEX IF NOT EXISTS idx_user_team_id ON "user"(team_id);
CREATE INDEX IF NOT EXISTS idx_user_assigned_pr_id ON "user"(assigned_pull_request_id);
CREATE INDEX IF NOT EXISTS idx_pull_request_author_id ON pull_request(author_id);
CREATE INDEX IF NOT EXISTS idx_pull_request_status ON pull_request(status);
CREATE INDEX IF NOT EXISTS idx_team_name ON team(team_name);
CREATE INDEX IF NOT EXISTS idx_user_username ON "user"(username);