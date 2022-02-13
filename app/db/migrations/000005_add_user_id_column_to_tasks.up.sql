ALTER TABLE tasks
ADD user_id CHAR(36) NOT NULL,
  ADD CONSTRAINT fk_tasks_tbl_user_id FOREIGN KEY (user_id) REFERENCES users(id);
