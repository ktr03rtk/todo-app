CREATE TABLE IF NOT EXISTS tasks(
  task_id INT PRIMARY KEY,
  name VARCHAR (50) NOT NULL,
  detail VARCHAR (300) NOT NULL,
  status TINYINT UNSIGNED NOT NULL,
  completion_date TIMESTAMP NULL DEFAULT NULL,
  deadline TIMESTAMP NULL DEFAULT NULL,
  notification_count TINYINT UNSIGNED NOT NULL,
  postponed_count TINYINT UNSIGNED NOT NULL
);
