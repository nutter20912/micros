-- Modify "posts" table
ALTER TABLE `posts` MODIFY COLUMN `user_id` longtext NOT NULL;
-- Create "comments" table
CREATE TABLE `comments` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  `user_id` longtext NOT NULL,
  `post_id` bigint unsigned NOT NULL,
  `content` varchar(500) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_posts_comments` (`post_id`),
  INDEX `idx_comments_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_posts_comments` FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
