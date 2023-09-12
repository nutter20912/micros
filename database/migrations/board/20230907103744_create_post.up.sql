-- Create "posts" table
CREATE TABLE `posts` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  `title` varchar(50) NOT NULL,
  `content` varchar(500) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `content` (`content`),
  INDEX `idx_posts_deleted_at` (`deleted_at`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
