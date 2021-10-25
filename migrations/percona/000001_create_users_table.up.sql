BEGIN;

CREATE TABLE users (
    row_id BIGINT AUTO_INCREMENT NOT NULL COMMENT 'record unique identifier',

    user_id VARCHAR(32) NOT NULL COMMENT 'user unique identifier',

    first_name VARCHAR(255) NOT NULL COMMENT 'user first name',
    last_name  VARCHAR(255) NOT NULL COMMENT 'user last name',

    created_at BIGINT NOT NULL COMMENT 'time when user was created',
    updated_at BIGINT COMMENT 'time when user was updated',

    PRIMARY KEY(row_id ASC),

    INDEX user_id_hash_idx USING HASH (user_id)
) COMMENT='stores information about user profile' ENGINE=InnoDB;

COMMIT;