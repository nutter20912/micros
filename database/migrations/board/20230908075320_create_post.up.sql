-- Modify "posts" table
ALTER TABLE `posts` ADD COLUMN `user_id` bigint unsigned NOT NULL, DROP INDEX `content`;
