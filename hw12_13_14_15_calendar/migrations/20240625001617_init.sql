-- +goose Up
-- +goose StatementBegin
create table events (
    event_id bigserial primary key,
    title varchar(300) not null check length(title) > 0,
    event_date timestamptz not null,
    duration interval not null check duration > 1,
    description varchar(1000),
    owner_user_id int8 not null check owner_user_id > 1,
    notify_time interval
);
comment on table event is "Событие в календаре";
comment on column event.event_id is "Уникальный идентификатор события";
comment on column event.title is "Заголовок";
comment on column event.event_date is "Дата и время события";
comment on column event.duration is "Длительность события";
comment on column event.description is "Описание события";
comment on column event.owner_user_id is "Код пользователя владельца события";
comment on column event.notify_time is "За сколько времени высылать уведомление?";

-- Поиск во всех случаях выполняется по дате события, время никогда не используется.
create index event_datetime_as_date_index on events ((event_date::date));
comment on index event_datetime_as_date_index is 'Индекс для поиска по дате события (без времени).';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table events;
-- +goose StatementEnd
