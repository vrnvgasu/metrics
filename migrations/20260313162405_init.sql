-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

CREATE TABLE IF NOT EXISTS metrics
(
    id    VARCHAR(250) NOT NULL,
    type  VARCHAR(50) NOT NULL,
    delta BIGINT,
    value DOUBLE PRECISION,
    hash  TEXT        NOT NULL DEFAULT '',
    PRIMARY KEY (id, type)
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

DROP TABLE IF EXISTS metrics;
