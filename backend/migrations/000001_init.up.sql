CREATE TABLE users (
    id         uuid PRIMARY KEY,
    google_sub text        NOT NULL UNIQUE,
    email      text        NOT NULL,
    name       text        NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE treatments (
    id         uuid PRIMARY KEY,
    user_id    uuid        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    treated_on date        NOT NULL,
    services   text[]      NOT NULL,
    salon_name text,
    memo       text,
    cost       integer,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT treatments_services_not_empty CHECK (cardinality(services) > 0),
    CONSTRAINT treatments_cost_non_negative CHECK (cost IS NULL OR cost >= 0)
);

CREATE INDEX treatments_user_id_treated_on_idx ON treatments (user_id, treated_on DESC);

CREATE TABLE share_links (
    id         uuid PRIMARY KEY,
    user_id    uuid        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token      text        NOT NULL UNIQUE,
    expires_at timestamptz NOT NULL,
    revoked_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX share_links_user_id_idx ON share_links (user_id);

CREATE TABLE sessions (
    id         uuid PRIMARY KEY,
    user_id    uuid        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token_hash text        NOT NULL UNIQUE,
    expires_at timestamptz NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX sessions_user_id_idx ON sessions (user_id);
CREATE INDEX sessions_expires_at_idx ON sessions (expires_at);
