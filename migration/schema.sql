CREATE TABLE users (
    id   BIGSERIAL PRIMARY KEY,
    name text      NOT NULL,
    authID text    NOT NULL,
    email text      NOT NULL
);

CREATE TABLE messages (
    id   BIGSERIAL PRIMARY KEY,
    from_id  BIGINT REFERENCES users(id) NOT NULL,
    from_authid  string NOT NULL,
    from_name text NOT NULL,
    message text NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
