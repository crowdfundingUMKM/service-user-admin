CREATE TABLE `users`
(
    `id` INT(11) NOT NULL AUTO_INCREMENT,
    `unix_id` CHAR(12) NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `email` VARCHAR(255),
    `phone` VARCHAR(20),
    `password_hash` varchar(255) DEFAULT NULL,
    `avatar_file_name` varchar(255) DEFAULT NULL,
    `status_account` varchar(255) DEFAULT NULL,
    `token` varchar(255) DEFAULT NULL,
    `ref_admin` varchar(255) DEFAULT NULL,
    `created_at` datetime DEFAULT NULL,
    `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- data
INSERT INTO `users` (`id`, `unix_id`,`name`, `email`, `phone`, `password_hash`, `avatar_file_name`, `status_account`, `token`, `ref_admin`, `created_at`, `updated_at`) VALUES
(1, '7d4aa4f2-90a', 'Ahmad Zaky', 'test@gmail.com', "82363152828", '$2a$04$6A5/psA4hCa0p0mLZQw4A.GKrkYDH3nTiim8lj9mYS18dmVi2FIvO', '', 'active', '','MASTER', '2023-03-15 22:56:25', '2023-03-15 22:56:25');

-- Indexes for table `users`
--
-- ALTER TABLE `users`
--   ADD PRIMARY KEY (`id`);

-- Remove token from table users
-- DELIMITER //

-- CREATE EVENT delete_expired_tokens
-- ON SCHEDULE EVERY 1 HOUR
-- DO
-- BEGIN
--     DELETE FROM users
--     WHERE token IS NOT NULL
--     AND created_at < NOW() - INTERVAL 2 DAY;
-- END //

-- DELIMITER ;

-- Backup database
-- SELECT *
-- INTO OUTFILE '/path/to/backup/users_backup.csv'
-- FIELDS TERMINATED BY ','
-- ENCLOSED BY '"'
-- LINES TERMINATED BY '\n'
-- FROM users;