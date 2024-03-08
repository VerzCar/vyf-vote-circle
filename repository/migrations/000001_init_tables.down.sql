BEGIN;

drop table user_options;

drop table rankings cascade;

drop table votes cascade;

drop table circle_voters cascade;

drop table circle_candidates cascade;

drop table circles cascade;

DROP TYPE circleStage;

DROP TYPE subscriptionPackage;

DROP TYPE placement;

DROP TYPE commitment;

COMMIT;