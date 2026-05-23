-- 001_create_users
CREATE TABLE IF NOT EXISTS users (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email       VARCHAR(255) UNIQUE NOT NULL,
    name        VARCHAR(100) NOT NULL,
    avatar_url  TEXT,
    provider    VARCHAR(20) NOT NULL,
    provider_id VARCHAR(255) NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    updated_at  TIMESTAMPTZ DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_provider ON users(provider, provider_id);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

-- 002_create_memory_spaces
CREATE TABLE IF NOT EXISTS memory_spaces (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id),
    title       VARCHAR(200) NOT NULL,
    description TEXT,
    cover_image TEXT,
    theme       VARCHAR(50) DEFAULT 'default',
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    updated_at  TIMESTAMPTZ DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_spaces_user_id ON memory_spaces(user_id);
CREATE INDEX IF NOT EXISTS idx_spaces_deleted_at ON memory_spaces(deleted_at);
CREATE INDEX IF NOT EXISTS idx_spaces_user_active ON memory_spaces(user_id, created_at DESC) WHERE deleted_at IS NULL;

-- 003_create_memories
CREATE TABLE IF NOT EXISTS memories (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id       UUID NOT NULL REFERENCES users(id),
    space_id      UUID NOT NULL REFERENCES memory_spaces(id),
    image_url     TEXT NOT NULL,
    thumbnail_url TEXT NOT NULL,
    title         VARCHAR(300) NOT NULL,
    content       TEXT,
    ai_content    TEXT,
    ai_model      VARCHAR(50),
    ai_prompt     TEXT,
    location      VARCHAR(200),
    memory_date   DATE NOT NULL,
    sort_order    INT DEFAULT 0,
    created_at    TIMESTAMPTZ DEFAULT NOW(),
    updated_at    TIMESTAMPTZ DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_memories_user_id ON memories(user_id);
CREATE INDEX IF NOT EXISTS idx_memories_space_id ON memories(space_id);
CREATE INDEX IF NOT EXISTS idx_memories_memory_date ON memories(space_id, memory_date);
CREATE INDEX IF NOT EXISTS idx_memories_deleted_at ON memories(deleted_at);
CREATE INDEX IF NOT EXISTS idx_memories_user_active ON memories(user_id, created_at DESC) WHERE deleted_at IS NULL;
