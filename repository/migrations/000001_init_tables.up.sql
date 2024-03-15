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

CREATE TYPE subscriptionPackage AS ENUM (
    'S',
    'M',
    'L'
    );

CREATE TYPE circleStage AS ENUM (
    'COLD',
    'HOT',
    'CLOSED'
    );

create table circles
(
    id           bigserial
        constraint circles_pkey
            primary key,
    name         varchar(40)                             not null,
    description  varchar(1200)                           not null,
    image_src    text                                    not null,
    private      boolean     default false               not null,
    active       boolean     default true                not null,
    stage        circleStage default 'COLD'::circleStage not null,
    created_from varchar(50)                             not null,
    valid_from   timestamp with time zone                not null,
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
    commitment   commitment default 'OPEN'::commitment not null,
    voted_for    varchar(50),
    circle_id    bigint
        constraint fk_circle_voters_circle
            references circles
            on delete restrict,
    circle_refer bigint
        constraint fk_circles_circle_voters
            references circles
            on delete restrict,
    created_at   timestamp with time zone,
    updated_at   timestamp with time zone,
    unique (voter, circle_id)
);

create table circle_candidates
(
    id           bigserial
        constraint circle_candidates_pkey
            primary key,
    candidate    varchar(50)                           not null,
    commitment   commitment default 'OPEN'::commitment not null,
    circle_id    bigint
        constraint fk_circle_candidates_circle
            references circles
            on delete restrict,
    circle_refer bigint
        constraint fk_circles_circle_candidates
            references circles
            on delete restrict,
    created_at   timestamp with time zone,
    updated_at   timestamp with time zone,
    unique (candidate, circle_id)
);

create table votes
(
    id              bigserial
        constraint votes_pkey
            primary key,
    voter_refer     bigint
        constraint fk_votes_voter_refer
            references circle_voters
            on delete restrict,
    candidate_refer bigint
        constraint fk_votes_candidate_refer
            references circle_candidates
            on delete restrict,
    circle_id       bigint
        constraint fk_votes_circle
            references circles
            on delete restrict,
    circle_refer    bigint
        constraint fk_circles_votes
            references circles
            on delete restrict,
    created_at      timestamp with time zone,
    updated_at      timestamp with time zone
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

create table user_options
(
    id                     bigserial
        constraint user_options_pkey
            primary key,
    identity_id            varchar(50)                                          not null,
    max_circles            int                                                  not null,
    max_voters             int                                                  not null,
    max_candidates         int                                                  not null,
    max_private_voters     int                                                  not null,
    max_private_candidates int                                                  not null,
    package                subscriptionPackage default 'S'::subscriptionPackage not null,
    created_at             timestamp with time zone,
    updated_at             timestamp with time zone
);

create index idx_user_options_id
    on user_options (id);

COMMIT;