-- Insert teams
INSERT INTO teams (name) VALUES
('backend'),
('frontend'),
('devops');

-- Insert users
INSERT INTO users (id, username, is_active, team_name) VALUES
-- Backend team users
('u1', 'alice', true, 'backend'),
('u5', 'eve', true, 'backend'),
('u6', 'frank', true, 'backend'),

-- Frontend team users
('u2', 'bob', true, 'frontend'),
('u7', 'grace', true, 'frontend'),
('u8', 'heidi', true, 'frontend'),

-- DevOps team users
('u4', 'dave', true, 'devops'),
('u9', 'ivan', true, 'devops'),
('u10', 'judy', true, 'devops'),

-- User without a team
('u3', 'carol', false, null);

-- Insert pull requests
INSERT INTO pull_requests (id, name, author_id, status, need_more_reviewers, created_at, merged_at) VALUES
('pr1', 'Add authentication', 'u1', 'open', true, now(), null),
('pr2', 'Fix UI bug', 'u2', 'merged', false, now() - interval '7 days', now() - interval '2 days'),
('pr3', 'Improve deployment script', 'u4', 'in_review', true, now() - interval '1 day', null);

-- Insert reviewers
INSERT INTO reviewers (user_id, pull_request_id) VALUES
('u2', 'pr1'),
('u3', 'pr1'),
('u1', 'pr3');
