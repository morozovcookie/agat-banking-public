BEGIN;

CREATE TABLE refresh_tokens (
    row_id BIGINT AUTO_INCREMENT NOT NULL COMMENT 'record unique identifier',

    token_id              VARCHAR(64) NOT NULL COMMENT 'token unique identifier',
    token_user_account_id VARCHAR(64) NOT NULL COMMENT 'user account unique identifier',
    token_issued_at       BIGINT NOT NULL COMMENT 'time when token was issued',
    token_expiration      BIGINT NOT NULL COMMENT 'time which after token will be expired (not for modification)',
    token_valid_until     BIGINT NOT NULL COMMENT 'time which after token will be expired (can be modified, use this
for finding expired tokens)',

    token_value VARCHAR(4096) NOT NULL COMMENT 'refresh token',

    created_at BIGINT NOT NULL COMMENT 'time when record was stored',
    updated_at BIGINT COMMENT 'time when record was updated',

    PRIMARY KEY (row_id DESC),

    INDEX token_id_hash_idx USING HASH (token_id) COMMENT 'use this for finding specific token, not value from
token_value column',
    INDEX token_valid_until_btree_idx USING BTREE (token_valid_until) COMMENT 'use this for finding expired tokens'
) COMMENT='stores user accounts refresh tokens' ENGINE=InnoDB;

COMMIT;
