-- HairHistoryMemo MVP schema (PostgreSQL)

CREATE TABLE users (
    id TEXT PRIMARY KEY,
    name TEXT,
    email TEXT,
    password_hash TEXT,
    last_login_at TIMESTAMPTZ,
    is_deactivated BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE hair_histories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id TEXT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    date DATE NOT NULL,
    services JSONB NOT NULL,
    salon_name TEXT NOT NULL DEFAULT '',
    stylist_name TEXT NOT NULL DEFAULT '',
    memo TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_hair_histories_user_date ON hair_histories (user_id, date DESC, created_at DESC);
