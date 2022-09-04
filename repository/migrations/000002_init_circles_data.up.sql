BEGIN;

INSERT INTO circles(name, description, image_src, created_from, created_at, updated_at)
VALUES ('Global',
        'One of the most impressive circles around the globe',
        '',
        'passoAvanti',
        now(),
        now());

COMMIT;