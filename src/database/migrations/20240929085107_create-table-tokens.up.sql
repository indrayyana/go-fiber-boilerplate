CREATE TABLE tokens(
    id              UUID            PRIMARY KEY DEFAULT uuid_generate_v4(),
    token           VARCHAR(255)    NOT NULL,
    user_id         UUID            NOT NULL,
    type            VARCHAR(255)    NOT NULL,
    expires         TIMESTAMP       DEFAULT CURRENT_TIMESTAMP  NOT NULL,
    created_at      TIMESTAMP       DEFAULT CURRENT_TIMESTAMP  NOT NULL,
    updated_at      TIMESTAMP       DEFAULT CURRENT_TIMESTAMP  NOT NULL,
    CONSTRAINT fk_user
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
