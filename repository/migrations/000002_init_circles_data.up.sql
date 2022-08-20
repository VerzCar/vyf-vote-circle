BEGIN;

INSERT INTO circles(name, description, image_src, created_from, created_at, updated_at)
VALUES ('Global', '', '', 'passoAvanti', now(), now());

COMMIT;