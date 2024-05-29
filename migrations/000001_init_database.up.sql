-- init tables

create table if not exists links
(
  id         bigserial primary key,
  created_at timestamp with time zone not null default now(),
  alias      text                     not null,
  long_url   text                     not null,
  clicks     int                      not null default 0
);
