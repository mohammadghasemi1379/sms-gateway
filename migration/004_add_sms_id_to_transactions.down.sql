ALTER TABLE transactions 
DROP FOREIGN KEY transactions_ibfk_2,
DROP INDEX idx_transactions_sms_id,
DROP COLUMN sms_id;
