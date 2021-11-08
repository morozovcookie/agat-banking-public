BEGIN;

ALTER TABLE user_accounts ADD INDEX account_id_hash_idx USING HASH (account_id);

COMMIT;
