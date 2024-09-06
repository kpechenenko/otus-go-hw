-- +goose Up
-- +goose StatementBegin
create table events (
    id uuid primary key,
    title varchar not null,
    date timestamp with time zone not null,
    duration interval not null,
    description text,
    owner_id int8 not null,
    notify_for interval
);
create index idx_events_date on events(date);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
drop table if exists events;
drop index if exists idx_events_date;
-- +goose StatementEnd
