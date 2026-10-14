CREATE TABLE IF NOT EXISTS comments(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    task_id UUID NOT NULL REFERENCES tasks(id)
ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id)
ON DELETE CASCADE,
    body TEXT NOT NULL CHECK(CHAR_LENGTH(body) BETWEEN 1
    AND 2000),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_comments_task_created
ON comments(
    task_id,
    created_at
);
CREATE INDEX idx_comments_user_id
ON comments(user_id)
