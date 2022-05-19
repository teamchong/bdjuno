CREATE TABLE upgrade_params
(
    one_row_id      BOOLEAN   NOT NULL DEFAULT TRUE PRIMARY KEY,
    binary_version  TEXT     NOT NULL,
    upgrade_info    TEXT      NOT NULL,
    upgrade_height  BIGINT    NOT NULL,
    upgrade_status  TEXT      NOT NULL,
    CHECK (one_row_id)
);