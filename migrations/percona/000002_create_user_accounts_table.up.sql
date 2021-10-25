BEGIN;

CREATE TABLE user_accounts (
    row_id BIGINT AUTO_INCREMENT NOT NULL COMMENT 'record unique identifier',

    account_id    VARCHAR(64)  NOT NULL COMMENT 'user account unique identifier',
    username      VARCHAR(255) NOT NULL COMMENT 'user login',
    email_address VARCHAR(255) NOT NULL COMMENT 'user email address',
    password_hash VARCHAR(255)  NOT NULL COMMENT 'user hashed password',

    user_id VARCHAR(64) NOT NULL COMMENT 'owner of account',

    created_at BIGINT NOT NULL COMMENT 'time when user account was created',
    updated_at BIGINT COMMENT 'time when user account was updated',

    PRIMARY KEY(row_id DESC),

    INDEX username_hash_idx      USING HASH (username),
    INDEX email_address_hash_idx USING HASH (email_address)
) COMMENT='stores user account information' ENGINE=InnoDB;

ALTER TABLE users MODIFY COLUMN user_id VARCHAR(64);

INSERT INTO users (user_id, first_name, last_name, created_at)
VALUES ('5g4agungfwzmmhrso7bdmk0rkz4kdfsrxebs855rvvejrp19qlc73zzuulh3fid2', 'admin', 'admin', 1635185552499);

INSERT INTO user_accounts (account_id, username, email_address, password_hash, user_id, created_at)
VALUES ('1s8cy82t3uiythh1on5vag79jx5737uii14xjbhdn5tn0a4g0psnnhytczqsihny', 'admin', 'admin@example.com',
        '8c6976e5b5410415bde908bd4dee15dfb167a9c873fc4bb8a81f6f2ab448a918',
        '5g4agungfwzmmhrso7bdmk0rkz4kdfsrxebs855rvvejrp19qlc73zzuulh3fid2', 1635185552499);

COMMIT;