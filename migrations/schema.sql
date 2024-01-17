CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY,
  name TEXT,
  email TEXT UNIQUE,
  token TEXT,
  password TEXT,
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

-- Passwords for users are as follows. Actual values that are inserted are encrypted.
-- abc123
-- def456
-- ghi789
-- jkl012
-- mno345
-- pqr678
-- stu901
-- vwx234
-- yz0123
-- 456xyz
-- 789abc
-- 012def
-- 345ghi
-- 678jkl
-- 901mno
INSERT INTO users (name, email, password, token, start, extra_vacation) VALUES
('John Doe', 'johndoe@example.com', '$2a$10$ccYRDLa1BZfGnexQARx6ceerQ2pTqj.Je/7GZAAItq2oDUVuqGrMS', 'abc123', '2022-01-01 09:00:00', 2),
('Jane Smith', 'janesmith@example.com', '$2a$10$z812onnie1BZBnR1PCRvm.ovfUMA.x8R1zulGFr9odjWZiwjzpsre', 'def456', '2021-12-15 14:30:00', 0),
('Michael Johnson', 'michaeljohnson@example.com', '$2a$10$d/cGBJeU9WVj52Stwfo8HOldAiIOVRDTWlz8IC3z8IPkC12E0okGG', 'ghi789', '2022-03-01 10:15:00', 1),
('Emily Davis', 'emilydavis@example.com', '$2a$10$UW9ZpoA3UktLrakpyG6CZOEDJB4AYktiIG0Dh12gi8okctEYsFn9O', 'jkl012', '2022-02-20 08:45:00', 3),
('Sarah Wilson', 'sarahwilson@example.com', '$2a$10$O9mxv.uHUCWpSAUYzg2PQe4DmKLH42hlGEceYJBfSKVCuqgHUaLM2', 'mno345', '2022-04-10 11:30:00', 1),
('David Thompson', 'davidthompson@example.com', '$2a$10$SXWtZr/Uam8utt06RU1BreIU2aixNk0.mXKxy0lvfpkSgvO61amzC', 'pqr678', '2022-05-05 09:15:00', 0),
('Olivia Garcia', 'oliviagarcia@example.com', '$2a$10$zhq8CiZ9XD7yuMWnkA58eOU8NigM6HOgsFqoRMFPXMDRHm5S8zQhC', 'stu901', '2022-03-18 14:00:00', 2),
('Jacob Martinez', 'jacobmartinez@example.com', '$2a$10$G76WtPchrflOC2WbxIkj2OL/z.w/VT45zIbN4VyZUrPOGIjqWTOPu', 'vwx234', '2022-01-10 10:30:00', 1),
('Emma Robinson', 'emmarobinson@example.com', '$2a$10$aKTRyHZmSszLkDf3MYDQNuW3ZG0CnKvUQ1jp2AyCJ5TDdBKS2KH3i', 'yz0123', '2022-06-25 09:45:00', 0),
('Noah Lee', 'noahlee@example.com', '$2a$10$GujvZ3WJ0zQoHJxCASXdkuXzb3GGAE.zeOD0UnTmgdtQYAFLe3oCW', '456xyz', '2022-07-15 08:00:00', 3),
('Ava Hernandez', 'avahernandez@example.com', '$2a$10$KHq/xHyva2aBLzlvsfQ24erSxcHzQcE2IbcwcE/xZenAmV9wy6qmK', '789abc', '2022-08-05 11:30:00', 1),
('William Clark', 'williamclark@example.com', '$2a$10$XNpcMV5LZXDXZHXJFOIklurHb0mUecojg8PO1Ab5323Eg8WU2Hdna', '012def', '2022-09-20 09:15:00', 0),
('Sophia Adams', 'sophiaadams@example.com', '$2a$10$fol.IizR9977H/YsEzaR2O5fz3qfo4nO.AbwiQmgwWfCD1yArLyEq', '345ghi', '2022-11-10 14:00:00', 2),
('James Baker', 'jamesbaker@example.com', '$2a$10$S/0cmP4Xtbk7fwxj1RMqbufCnC2bWwg9RZfN4sUkt8m3m7e1eO6Ha', '678jkl', '2022-12-12 10:30:00', 1),
('Mia Mitchell', 'miamitchell@example.com', '$2a$10$x6UmBUEtKPoZeRbtEBoBSel9uDu/bA.fpCm2bHJOe2rYb.wJPEw6G', '901mno', '2022-10-18 09:45:00', 0);


INSERT INTO leaves (user_id, start, end, type, approved) VALUES
(1, '2023-07-01 09:00:00', '2023-07-05 18:00:00', 'vacation', false),
(2, '2023-07-15 12:00:00', '2023-07-17 16:30:00', 'sick', false),
(3, '2023-07-10 08:30:00', '2023-07-12 17:00:00', 'vacation', false),
(4, '2023-07-20 14:00:00', '2023-07-22 12:30:00', 'dayoff', false),
(5, '2023-07-10 09:00:00', '2023-07-11 18:00:00', 'vacation', false),
(6, '2023-07-05 13:00:00', '2023-07-07 15:30:00', 'sick', false),
(7, '2023-07-15 08:30:00', '2023-07-18 17:00:00', 'vacation', false),
(8, '2023-07-10 10:00:00', '2023-07-11 14:30:00', 'dayoff', false),
(9, '2023-07-20 11:30:00', '2023-07-23 12:00:00', 'vacation', false),
(3, '2023-10-05 08:00:00', '2023-07-08 18:00:00', 'sick', false),
(2, '2022-11-15 12:30:00', '2022-11-17 16:00:00', 'vacation', false),
(4, '2022-12-10 09:30:00', '2022-12-12 13:30:00', 'dayoff', false),
(1, '2023-01-05 14:00:00', '2023-01-06 17:30:00', 'vacation', false),
(3, '2023-02-10 08:30:00', '2023-02-12 16:00:00', 'dayoff', false),
(2, '2023-03-15 09:00:00', '2023-03-17 15:00:00', 'vacation', false);
