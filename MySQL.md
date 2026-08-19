# MySQL
To connect MySQL: 
```sh
brew services start mysql
mysql -u root
```
To stop MySQL: `brew services stop mysql`

## Setup database
Create `snippetbox` database and start using it:
```sh
CREATE DATABASE snippetbox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE snippetbox
```
## Create a new `snippets` table
```sh
CREATE TABLE snippets (
  id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
  title VARCHAR(100) NOT NULL,
  content TEXT NOT NULL,
  created DATETIME NOT NULL,
  expires DATETIME NOT NULL
);
```
Convert column `created` to INDEX:
```sh
CREATE INDEX idx_snippets_created ON snippets(created);
```

## Add data to table
Add data to `snippets` table:
```sh
INSERT INTO snippets (title, content, created, expires) VALUES (
  'An old silent pond',
  'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n– Matsuo Bashō',
  UTC_TIMESTAMP(),
  DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)
);
```
```sh
INSERT INTO snippets (title, content, created, expires) VALUES (
'Over the wintry forest',
'Over the wintry\nforest, winds howl in rage\nwith no leaves to blow.\n\n– Natsume Soseki',
UTC_TIMESTAMP(),
DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)
);
```
```sh
INSERT INTO snippets (title, content, created, expires) VALUES (
'First autumn morning',
'First autumn morning\nthe mirror I stare into\nshows my father''s face.\n\n– Murakami Kijo',
UTC_TIMESTAMP(),
DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)
);
```

## Add User with certain privileges
Create User named `abc` with password `1234`
```sh
CREATE USER 'abc'@'localhost' IDENTIFIED BY '1234';
```
Grant privileges:
```sh
GRANT SELECT, INSERT, UPDATE, DELETE ON snippetbox.* TO 'abc'@'localhost';
```
Logout from current User: `exit`.\
Login to database `snippetbox` from that User: `mysql -D snippetbox -u abc -p` => type password `1234`.\
Read data from a table: `SELECT id, title, expires FROM snippets` 

## Session on Server side
Create `sessions` table
```sh
USE snippetbox;

CREATE TABLE sessions (
  token CHAR(43) PRIMARY KEY,
  data BLOB NOT NULL,
  expiry TIMESTAMP(6) NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions (expiry);
```

## Create `users` table
```sh
USE snippetbox;

CREATE TABLE users (
  id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL,
  hashed_password CHAR(60) NOT NULL,
  created DATETIME NOT NULL
);

ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);
```
Test credentials for login: `abc@gmail.com` and `12345678`

## Delete an User from table
```sh
USE snippetbox;
DELETE FROM users WHERE email="abc@gmail.com";
```