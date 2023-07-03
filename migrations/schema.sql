CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY,
  name TEXT,
  email TEXT UNIQUE,
  token TEXT,
  start DATETIME,
  extra_vacation INTEGER
);

CREATE TABLE IF NOT EXISTS leaves (
  id INTEGER PRIMARY KEY,
  user_id INTEGER,
  start DATETIME,
  end DATETIME,
  type TEXT,
  approved BOOLEAN,
  FOREIGN KEY(user_id) REFERENCES users(id)
);

INSERT INTO users (name, email, token, start, extra_vacation) VALUES
('John Doe', 'johndoe@example.com', 'abc123', '2022-01-01 09:00:00', 2),
('Jane Smith', 'janesmith@example.com', 'def456', '2021-12-15 14:30:00', 0),
('Michael Johnson', 'michaeljohnson@example.com', 'ghi789', '2022-03-01 10:15:00', 1),
('Emily Davis', 'emilydavis@example.com', 'jkl012', '2022-02-20 08:45:00', 3),
('Sarah Wilson', 'sarahwilson@example.com', 'mno345', '2022-04-10 11:30:00', 1),
('David Thompson', 'davidthompson@example.com', 'pqr678', '2022-05-05 09:15:00', 0),
('Olivia Garcia', 'oliviagarcia@example.com', 'stu901', '2022-03-18 14:00:00', 2),
('Jacob Martinez', 'jacobmartinez@example.com', 'vwx234', '2022-01-10 10:30:00', 1),
('Emma Robinson', 'emmarobinson@example.com', 'yz0123', '2022-06-25 09:45:00', 0),
('Noah Lee', 'noahlee@example.com', '456xyz', '2022-07-15 08:00:00', 3),
('Ava Hernandez', 'avahernandez@example.com', '789abc', '2022-08-05 11:30:00', 1),
('William Clark', 'williamclark@example.com', '012def', '2022-09-20 09:15:00', 0),
('Sophia Adams', 'sophiaadams@example.com', '345ghi', '2022-11-10 14:00:00', 2),
('James Baker', 'jamesbaker@example.com', '678jkl', '2022-12-12 10:30:00', 1),
('Mia Mitchell', 'miamitchell@example.com', '901mno', '2022-10-18 09:45:00', 0);


INSERT INTO leaves (user_id, start, end, type, approved) VALUES
(1, '2022-01-01 09:00:00', '2022-01-05 18:00:00', 'vacation', false),
(2, '2022-02-15 12:00:00', '2022-02-17 16:30:00', 'sick', false),
(3, '2022-03-10 08:30:00', '2022-03-12 17:00:00', 'vacation', false),
(4, '2022-04-20 14:00:00', '2022-04-22 12:30:00', 'dayoff', false),
(1, '2022-05-10 09:00:00', '2022-05-11 18:00:00', 'vacation', false),
(3, '2022-06-05 13:00:00', '2022-06-07 15:30:00', 'vacation', false),
(2, '2022-07-15 08:30:00', '2022-07-18 17:00:00', 'vacation', false),
(4, '2022-08-10 10:00:00', '2022-08-11 14:30:00', 'vacation', false),
(1, '2022-09-20 11:30:00', '2022-09-23 12:00:00', 'vacation', false),
(3, '2022-10-05 08:00:00', '2022-10-08 18:00:00', 'sick', false),
(2, '2022-11-15 12:30:00', '2022-11-17 16:00:00', 'vacation', false),
(4, '2022-12-10 09:30:00', '2022-12-12 13:30:00', 'dayoff', false),
(1, '2023-01-05 14:00:00', '2023-01-06 17:30:00', 'vacation', false),
(3, '2023-02-10 08:30:00', '2023-02-12 16:00:00', 'dayoff', false),
(2, '2023-03-15 09:00:00', '2023-03-17 15:00:00', 'vacation', false);
