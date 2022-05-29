BEGIN;

INSERT INTO circles(name, created_from, created_at, updated_at)
VALUES ('Global', 'af', now(), now());

INSERT INTO circle_voters(voter, committed, circle_id, circle_refer, created_at, updated_at)
VALUES ('af', true, 1, 1, now(), now());

COMMIT;