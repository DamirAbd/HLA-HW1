CREATE TABLE IF NOT EXISTS users (
  userID serial NOT NULL,
  ID VARCHAR(255) NOT NULL,
  FirstName VARCHAR(255) NOT NULL,
  SecondName VARCHAR(255) NOT NULL,
  BirthDate DATE NOT NULL,
  Biography VARCHAR(255) NULL,
  City VARCHAR(255) NULL,
  Password VARCHAR(255) NOT NULL,
  CreatedAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  
  PRIMARY KEY (userID)
);

CREATE TABLE IF NOT EXISTS friends (
  fr_id serial NOT NULL,
  user_id VARCHAR(255) REFERENCES users (userID),
  friend_id VARCHAR(255) REFERENCES users (userID),

  PRIMARY KEY (fr_id)

);

CREATE TABLE IF NOT EXISTS posts (
  post_id VARCHAR(255) NOT NULL,
  author_id character varying(255) REFERENCES users (ID),
  post TEXT NULL,

  PRIMARY KEY (post_id)

);