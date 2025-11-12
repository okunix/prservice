CREATE TABLE IF NOT EXISTS teams (
    name VARCHAR PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS users (
    id VARCHAR PRIMARY KEY,
    name VARCHAR NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    teamName VARCHAR NULL DEFAULT NULL,
    FOREIGN KEY (teamName) REFERENCES teams(name) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS pull_requests (
    id VARCHAR PRIMARY KEY,
    name VARCHAR NOT NULL,
    author_id VARCHAR NOT NULL,
    status VARCHAR NOT NULL,
    need_more_reviewers BOOLEAN DEFAULT false,
    FOREIGN KEY (author_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS reviewers (
    user_id VARCHAR, 
    pull_request_id VARCHAR,
    PRIMARY KEY (user_id, pull_request_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (pull_request_id) REFERENCES pull_requests(id) ON DELETE CASCADE
);
