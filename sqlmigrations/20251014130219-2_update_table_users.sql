-- +migrate Up
-- Tạo ENUM type cho user status
CREATE TYPE user_status AS ENUM ('pending', 'active', 'blocked', 'banned');

-- Drop constraint cũ
ALTER TABLE users
    DROP CONSTRAINT IF EXISTS users_status_check;

-- Drop DEFAULT value cũ (TEXT)
ALTER TABLE users
    ALTER COLUMN status DROP DEFAULT;

-- Chuyển column type từ TEXT sang ENUM
ALTER TABLE users
    ALTER COLUMN status TYPE user_status USING status::user_status;

-- Set DEFAULT mới (ENUM)
ALTER TABLE users
    ALTER COLUMN status SET DEFAULT 'pending'::user_status;

-- +migrate Down
-- Drop DEFAULT ENUM
ALTER TABLE users
    ALTER COLUMN status DROP DEFAULT;

-- Chuyển về TEXT
ALTER TABLE users
    ALTER COLUMN status TYPE TEXT;

-- Set DEFAULT TEXT
ALTER TABLE users
    ALTER COLUMN status SET DEFAULT 'active';

-- Add lại constraint cũ
ALTER TABLE users
    ADD CONSTRAINT users_status_check CHECK (status IN ('active', 'blocked', 'banned'));

-- Xóa ENUM type
DROP TYPE IF EXISTS user_status;