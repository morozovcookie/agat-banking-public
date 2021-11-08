BEGIN;

ALTER TABLE user_accounts DROP INDEX account_id_hash_idx;

COMMIT;
