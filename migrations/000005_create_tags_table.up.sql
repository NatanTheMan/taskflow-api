CREATE TABLE if NOT EXISTS tags(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id)
ON DELETE cascade,
    name TEXT NOT NULL CHECK(CHAR_LENGTH(name) BETWEEN 1
    AND 50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT uq_tags_user_name UNIQUE(
        user_id,
        name
    )
);
