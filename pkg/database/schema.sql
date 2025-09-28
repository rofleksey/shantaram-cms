CREATE TABLE IF NOT EXISTS orders
(
    id             UUID PRIMARY KEY,
    index          BIGSERIAL    NOT NULL,
    table_id       VARCHAR(255),
    created        TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated        TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    status         VARCHAR(64)  NOT NULL,
    client_name    VARCHAR(255) NOT NULL,
    client_phone   VARCHAR(64)  NOT NULL,
    client_comment TEXT,
    seen           BOOLEAN      NOT NULL DEFAULT false,
    items          JSONB        NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_orders_created ON orders (index DESC);

CREATE TABLE IF NOT EXISTS menu
(
    id      VARCHAR(64) PRIMARY KEY,
    title   VARCHAR(255) NOT NULL,
    created TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS product_groups
(
    id      UUID PRIMARY KEY,
    menu_id VARCHAR(64)  NOT NULL REFERENCES menu (id) ON DELETE CASCADE,
    index   INTEGER      NOT NULL CHECK (index >= 0),
    title   VARCHAR(255) NOT NULL,
    created TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT product_groups_order UNIQUE (menu_id, index) DEFERRABLE INITIALLY DEFERRED
);

CREATE TABLE IF NOT EXISTS products
(
    id          UUID PRIMARY KEY,
    group_id    UUID             NOT NULL REFERENCES product_groups (id) ON DELETE CASCADE,
    index       INTEGER          NOT NULL CHECK (index >= 0),
    title       VARCHAR(255)     NOT NULL,
    description TEXT             NOT NULL,
    price       DOUBLE PRECISION NOT NULL CHECK (price >= 0.0),
    available   BOOLEAN          NOT NULL DEFAULT true,
    created     TIMESTAMP        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated     TIMESTAMP        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT products_order UNIQUE (group_id, index) DEFERRABLE INITIALLY DEFERRED
);

CREATE TABLE IF NOT EXISTS migration
(
    id      VARCHAR(255) PRIMARY KEY,
    applied TIMESTAMP NOT NULL
);
