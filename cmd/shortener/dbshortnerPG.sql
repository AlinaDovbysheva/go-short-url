-- Database: DB_shortner

-- DROP DATABASE IF EXISTS "DB_shortner";
--create extension pgcrypto;

CREATE TABLE IF NOT EXISTS users(
                      id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
                      user_id uuid NOT NULL,
                      key text NOT NULL DEFAULT gen_salt('md5'),
                      uname VARCHAR ( 255 ) NULL,
                      last_login TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS url (
                     id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
                     url VARCHAR ( 255 ) UNIQUE NOT NULL,
                     url_short VARCHAR ( 255 ) UNIQUE NOT NULL,
                     created_on TIMESTAMP NOT NULL  DEFAULT NOW()
);


CREATE TABLE IF NOT EXISTS users_url (
                           url_id INT NOT NULL,
                           user_id INT NOT NULL,
                           FOREIGN KEY(url_id)
                               REFERENCES url(id)
                               ON DELETE SET NULL,

                           FOREIGN KEY (user_id)
                               REFERENCES users (id)
);



