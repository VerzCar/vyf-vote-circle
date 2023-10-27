BEGIN;

CREATE
    EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE commitment AS ENUM (
    'OPEN',
    'COMMITTED',
    'REJECTED'
    );

CREATE TYPE placement AS ENUM (
    'NEUTRAL',
    'ASCENDING',
    'DESCENDING'
    );

create table circles
(
    id           bigserial
        constraint circles_pkey
            primary key,
    name         varchar(40)           not null,
    description  varchar(1200)         not null,
    image_src    text                  not null,
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
    voter        varchar(50)                           not null,
    voted_for    varchar(50),
    voted_from   varchar(50),
    commitment   commitment default 'OPEN'::commitment not null,
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
    id            bigserial
        constraint votes_pkey
            primary key,
    voter_refer   bigint
        constraint fk_votes_voter_refer
            references circle_voters
            on delete restrict,
    elected_refer bigint
        constraint fk_votes_elected_refer
            references circle_voters
            on delete restrict,
    circle_id     bigint
        constraint fk_votes_circle
            references circles
            on delete restrict,
    circle_refer  bigint
        constraint fk_circles_votes
            references circles
            on delete restrict,
    created_at    timestamp with time zone,
    updated_at    timestamp with time zone
);

create index idx_votes_id
    on votes (id);

create table rankings
(
    id          bigserial
        constraint rankings_pkey
            primary key,
    identity_id varchar(50)                            not null,
    number      int                                    not null,
    votes       bigint    default 0                    not null,
    placement   placement default 'NEUTRAL'::placement not null,
    circle_id   bigint
        constraint fk_rankings_circle
            references circles
            on delete restrict,
    created_at  timestamp with time zone,
    updated_at  timestamp with time zone
);

create index idx_rankings_id
    on rankings (id);

COMMIT;