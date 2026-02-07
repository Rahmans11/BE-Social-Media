-- CREATE TABLE "accounts" (
--   "id" INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
--   "email" VARCHAR(255) NOT NULL UNIQUE,
--   "password" VARCHAR(255) NOT NULL,
--   "role" VARCHAR(255) DEFAULT 'USER',
--   "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
--   "update_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
-- );

-- ALTER TABLE accounts
-- ADD COLUMN role VARCHAR(255) DEFAULT 'USER';

-- ALTER TABLE accounts
-- ALTER COLUMN role SET DEFAULT 'USER';
-- DELETE FROM accounts

TABLE accounts;

-- ALTER TABLE account ADD CONSTRAINT unique_email UNIQUE (email);

-- ALTER TABLE account ALTER COLUMN email SET NOT NULL;
-- ALTER TABLE account ALTER COLUMN password SET NOT NULL;

-- CREATE TABLE "users" (
--   "id" INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
--   "account_id" INT NOT NULL,
--   "first_name" VARCHAR(255),
--   "last_name" VARCHAR(255),
--   "email" VARCHAR(255) NOT NULL,
--   "phone_number" VARCHAR(255),
--   "avatar" VARCHAR(255),
--   "bio" TEXT,
--   "last_login" TIMESTAMPTZ,
--   "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
--   "update_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
--   FOREIGN KEY (account_id) REFERENCES accounts(id)
-- );

-- ALTER TABLE users DROP COLUMN role;

-- ALTER TABLE users ALTER COLUMN account_id SET NOT NULL;
-- ALTER TABLE users ALTER COLUMN email SET NOT NULL;

TABLE users;

--DELETE FROM users;

-- CREATE TABLE "posts" (
--   "id" INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
--   "user_id" INT NOT NULL,
--   "caption" TEXT NOT NULL,
--   "image" VARCHAR(255),
--   "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
--   "update_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
--   FOREIGN KEY (user_id) REFERENCES users(id)
-- );

-- ALTER TABLE posts ALTER COLUMN account_id SET NOT NULL;
-- ALTER TABLE posts ALTER COLUMN caption SET NOT NULL;
--ALTER TABLE posts RENAME COLUMN account_id TO user_id;

TABLE posts;

-- CREATE TABLE "posts_reactions" (
--   "post_id" INT NOT NULL,
--   "liked_id" INT NOT NULL,
--   "comment" TEXT,
--   "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
--   "update_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
--   FOREIGN KEY (post_id) REFERENCES posts(id),
--   FOREIGN KEY (liked_id) REFERENCES users(id) 
-- );

TABLE posts_reactions;

-- CREATE TABLE "follows" (
--   "followed_id" INT NOT NULL,
--   "follower_id" INT NOT NULL,
--   "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
--   "update_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
--   FOREIGN KEY (follower_id) REFERENCES users(id),
--   FOREIGN KEY (followed_id) REFERENCES users(id) 
-- );

TABLE follows;

-- ALTER TABLE follows
-- ADD CONSTRAINT fk_follower
-- FOREIGN KEY (follower_id)
-- REFERENCES accounts(id) 

-- ALTER TABLE follows
-- ADD CONSTRAINT fk_followed
-- FOREIGN KEY (followed_id)
-- REFERENCES accounts(id) 