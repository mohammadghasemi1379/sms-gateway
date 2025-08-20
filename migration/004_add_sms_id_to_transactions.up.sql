ALTER TABLE transactions 
ADD COLUMN sms_id INT UNSIGNED NULL,
ADD FOREIGN KEY (sms_id) REFERENCES sms(id) ON DELETE SET NULL,
ADD INDEX idx_transactions_sms_id (sms_id);
