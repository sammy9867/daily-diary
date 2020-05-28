DROP DATABASE IF EXISTS diary_db;
CREATE DATABASE diary_db;
USE diary_db;

CREATE TABLE users (
    id BIGINT AUTO_INCREMENT, 
    username VARCHAR(255) UNIQUE, 
    email VARCHAR(100), 
    password VARCHAR(100),
    created_at DATETIME, 
    updated_at DATETIME,
    PRIMARY KEY (id)
);

CREATE TABLE entries (
    id BIGINT AUTO_INCREMENT, 
    title VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    owner_id BIGINT NOT NULL,
    created_at DATETIME, 
    updated_at DATETIME,
    PRIMARY KEY (id),
    FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE entry_images (
    id BIGINT AUTO_INCREMENT, 
    url VARCHAR(255) NOT NULL,
    entry_id BIGINT NOT NULL,
    created_at DATETIME, 
    updated_at DATETIME,
    PRIMARY KEY (id),
    FOREIGN KEY (entry_id) REFERENCES entries(id) ON DELETE CASCADE
);

DROP PROCEDURE IF EXISTS GetAllEntryImagesOfEntry;
DELIMITER //
CREATE PROCEDURE GetAllEntryImagesOfEntry(IN entry_id BIGINT)
BEGIN
    SELECT * FROM entry_images
    WHERE entry_images.entry_id = entry_id;
END //
DELIMITER ;


