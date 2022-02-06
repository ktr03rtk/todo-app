ALTER TABLE sessions
ADD CONSTRAINT fk_sessions_tbl_user_id FOREIGN KEY (user_id) REFERENCES users(id);
