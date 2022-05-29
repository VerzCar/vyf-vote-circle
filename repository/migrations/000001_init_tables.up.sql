BEGIN;

CREATE
    EXTENSION IF NOT EXISTS "uuid-ossp";

create table circles
(
    id           bigserial
        constraint circles_pkey
            primary key,
    name         varchar(40)           not null,
    private      boolean default false not null,
    active       boolean default true  not null,
    created_from varchar(50)           not null,
    valid_until  timestamp with time zone,
    created_at   timestamp with time zone,
    updated_at   timestamp with time zone
);

create index idx_circles_id
    on circles (id);

create table circle_voters
(
    id           bigserial
        constraint circle_voters_pkey
            primary key,
    voter        varchar(50)           not null,
    committed    boolean default false not null,
    rejected     boolean default false not null,
    circle_id    bigint
        constraint fk_circle_voters_circle
            references circles
            on delete restrict,
    circle_refer bigint
        constraint fk_circles_circle_voters
            references circles
            on delete restrict,
    created_at   timestamp with time zone,
    updated_at   timestamp with time zone
);

create index idx_circle_voters_id
    on circle_voters (id);

create table votes
(
    id           bigserial
        constraint votes_pkey
            primary key,
    voter        varchar(50) not null,
    elected      varchar(50) not null,
    circle_id    bigint
        constraint fk_votes_circle
            references circles
            on delete restrict,
    circle_refer bigint
        constraint fk_circles_votes
            references circles
            on delete restrict,
    created_at   timestamp with time zone,
    updated_at   timestamp with time zone
);

create index idx_votes_id
    on votes (id);

COMMIT;