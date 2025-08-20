ALTER TABLE users ADD COLUMN phone_number VARCHAR(20) NOT NULL AFTER name;
CREATE UNIQUE INDEX idx_users_phone_number ON users(phone_number);
