-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

CREATE TABLE IF NOT EXISTS metrics
(
    id        BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    metric_id VARCHAR(250) NOT NULL,
    type      VARCHAR(50)  NOT NULL,
    delta     BIGINT,
    value     DOUBLE PRECISION,
    hash      VARCHAR(500) NOT NULL DEFAULT ''
);
CREATE UNIQUE INDEX IF NOT EXISTS id_unique_metrics_metric_id_type
    ON metrics (metric_id, type);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

DROP TABLE IF EXISTS metrics;
