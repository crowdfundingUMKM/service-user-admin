CREATE TABLE `users`
(
    `id` INT(11) NOT NULL AUTO_INCREMENT,
    `unix_id` CHAR(12) NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `email` VARCHAR(255),
    `phone` VARCHAR(20),
    `password_hash` VARCHAR(255),
    `avatar_file_name` VARCHAR(255),
    `status_account` VARCHAR(255),
    `token` VARCHAR(255),
    `ref_admin` VARCHAR(12) DEFAULT NULL,
    `update_id_admin` CHAR(12) DEFAULT NULL,
    `update_at_admin` DATETIME, DEFAULT NULL,
    `created_at` DATETIME DEFAULT NULL,
    `updated_at` DATETIME DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


CREATE TABLE `notif_admins`
(
    `id` INT(11) NOT NULL AUTO_INCREMENT,
    `user_admin_id` CHAR(12) NOT NULL,
    `title` VARCHAR(255) NOT NULL,
    `description` TEXT,
    `to_user` VARCHAR(11),
    `document` VARCHAR(255),
    `status_notif` TINYINT(1),
    `created_at` DATETIME DEFAULT NULL,
    `updated_at` DATETIME DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

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