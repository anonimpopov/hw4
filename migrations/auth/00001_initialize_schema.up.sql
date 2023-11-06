CREATE TABLE IF NOT EXISTS users (
   login         varchar(120) PRIMARY KEY,
   password_hash varchar(72) NOT NULL,
   email         varchar(300) UNIQUE NOT NULL
);