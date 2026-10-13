CREATE TABLE tasks(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id)
ON DELETE CASCADE,
    title TEXT NOT NULL check(CHAR_LENGTH(title) BETWEEN 3 and 120),
    description TEXT NOT NULL DEFAULT '' check(CHAR_LENGTH(description) <= 2000),
    status TEXT NOT NULL DEFAULT 'todo'CHECK(status IN('todo',
    'doing',
    'done')),
    priority TEXT NOT NULL DEFAULT 'medium'CHECK(priority IN('low',
    'medium',
    'high')),
    due_date timestamptz,
    completed_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()

    constraint chk_tasks_completed check((status='done') =(completed_at is not null))
);
CREATE index idx_tasks_project_status
ON tasks(project_id,
status);
CREATE index idx_tasks_project_due_date
ON tasks(project_id,
due_date);
