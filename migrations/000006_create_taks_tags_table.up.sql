CREATE TABLE IF NOT EXISTS task_tags (
task_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
tag_id  UUID NOT NULL REFERENCES tags(id)  ON DELETE CASCADE,
PRIMARY KEY (task_id, tag_id)
);

CREATE INDEX idx_task_tags_tag_id ON task_tags (tag_id);
