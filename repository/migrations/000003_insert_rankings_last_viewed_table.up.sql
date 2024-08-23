BEGIN;

create table rankings_last_viewed
(
    id          bigserial
        constraint rankings_last_viewed_pkey
            primary key,
    identity_id varchar(50) not null,
    circle_id   bigint
        constraint fk_rankings_last_viewed_circle
            references circles
            on delete restrict,
    created_at  timestamp with time zone,
    updated_at  timestamp with time zone
);

COMMIT;
