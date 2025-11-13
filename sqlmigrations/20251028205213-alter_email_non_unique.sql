
-- +migrate Up
/*
Xoá constraint tên users_email_key.
Khi PostgreSQL tạo UNIQUE constraint từ UNIQUE inline, tên mặc định thường là:
<table>_<column>_key
nên ở đây là users_email_key.
*/
ALTER TABLE users
    DROP CONSTRAINT IF EXISTS users_email_key;
/*
DROP field 
DROP NOT NULL
*/
ALTER TABLE users    
    ALTER COLUMN email DROP NOT NULL;
-- +migrate Down
ALTER TABLE users
    ALTER COLUMN email SET NOT NULL;

ALTER TABLE users
    ADD CONSTRAINT users_email_key UNIQUE(email);