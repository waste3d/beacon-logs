-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

create table page_views (
    url text primary key,
    view_count BIGINT NOT NULL DEFAULT 1,
    last_seen TIMESTAMPTZ NOT NULL
);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
drop table if exists page_views