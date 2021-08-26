-- UP
-- Threads
CREATE TABLE threads
(
  id SERIAL NOT NULL,
  is_deleted BOOLEAN NOT NULL,
  created_at TIMESTAMP DEFAULT NOW() NOT NULL,
  updated_at TIMESTAMP DEFAULT NOW() NOT NULL,
  PRIMARY KEY (id)
);

-- Posts
CREATE TABLE posts
(
  id SERIAL NOT NULL,
  thread_id INTEGER NOT NULL,
  title VARCHAR NOT NULL,
  text TEXT NOT NULL,
  is_sage BOOLEAN NOT NULL,
	created_at TIMESTAMP DEFAULT NOW() NOT NULL,
  PRIMARY KEY (id),
  FOREIGN KEY (thread_id) REFERENCES threads (id)
);


-- Files
CREATE TABLE files
(
  id SERIAL PRIMARY KEY,
  post_id INTEGER NOT NULL,
  created_at TIMESTAMP DEFAULT NOW() NOT NULL,
  name VARCHAR NOT NULL,
  ext VARCHAR NOT NULL,
  size INTEGER NOT NULL,
  FOREIGN KEY (post_id) REFERENCES posts (id)
);



-- DOWN
DROP TABLE files;
DROP TABLE posts;
DROP TABLE threads;