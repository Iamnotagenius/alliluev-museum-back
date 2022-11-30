DROP TABLE IF EXISTS exhibits;
DROP TABLE IF EXISTS rooms;

CREATE TABLE rooms (
    id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
    name CHAR(64) NOT NULL,
    pictures VARCHAR(128) COMMENT 'comma-separated list of paths'
);

CREATE TABLE exhibits (
    id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
    name CHAR(64) NOT NULL,
    pictures VARCHAR(128) COMMENT 'comma-separated list of paths',
    description TEXT,
    room INT,

    FOREIGN KEY (room) REFERENCES rooms(id)
        ON UPDATE CASCADE ON DELETE SET NULL
);
