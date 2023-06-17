CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY,
  name TEXT,
  email TEXT UNIQUE,
  token TEXT
);


INSERT INTO users (name, email, token) VALUES ('John Doe', 'john.doe@example.com', 'token123');
INSERT INTO users (name, email, token) VALUES ('Jane Smith', 'jane.smith@example.com', 'token456');
INSERT INTO users (name, email, token) VALUES ('Michael Johnson', 'michael.johnson@example.com', 'token789');
INSERT INTO users (name, email, token) VALUES ('Emily Davis', 'emily.davis@example.com', 'tokenabc');
INSERT INTO users (name, email, token) VALUES ('Daniel Wilson', 'daniel.wilson@example.com', 'tokendef');
INSERT INTO users (name, email, token) VALUES ('Sarah Brown', 'sarah.brown@example.com', 'tokenxyz');
INSERT INTO users (name, email, token) VALUES ('David Lee', 'david.lee@example.com', 'token123');
INSERT INTO users (name, email, token) VALUES ('Olivia Taylor', 'olivia.taylor@example.com', 'token456');
INSERT INTO users (name, email, token) VALUES ('Ethan Anderson', 'ethan.anderson@example.com', 'token789');
INSERT INTO users (name, email, token) VALUES ('Sophia Martinez', 'sophia.martinez@example.com', 'tokenabc');
INSERT INTO users (name, email, token) VALUES ('William Thompson', 'william.thompson@example.com', 'tokendef');
INSERT INTO users (name, email, token) VALUES ('Ava Garcia', 'ava.garcia@example.com', 'tokenxyz');
INSERT INTO users (name, email, token) VALUES ('James Rodriguez', 'james.rodriguez@example.com', 'token123');
INSERT INTO users (name, email, token) VALUES ('Mia Hernandez', 'mia.hernandez@example.com', 'token456');
INSERT INTO users (name, email, token) VALUES ('Benjamin Nelson', 'benjamin.nelson@example.com', 'token789');
INSERT INTO users (name, email, token) VALUES ('Abigail Johnson', 'abigail.johnson@example.com', 'tokenabc');
INSERT INTO users (name, email, token) VALUES ('Alexander Thomas', 'alexander.thomas@example.com', 'tokendef');
INSERT INTO users (name, email, token) VALUES ('Sofia Roberts', 'sofia.roberts@example.com', 'tokenxyz');
INSERT INTO users (name, email, token) VALUES ('Joseph Smith', 'joseph.smith@example.com', 'token123');
INSERT INTO users (name, email, token) VALUES ('Harper Williams', 'harper.williams@example.com', 'token456');
INSERT INTO users (name, email, token) VALUES ('Daniel Brown', 'daniel.brown@example.com', 'token789');
INSERT INTO users (name, email, token) VALUES ('Madison Johnson', 'madison.johnson@example.com', 'tokenabc');
INSERT INTO users (name, email, token) VALUES ('Sebastian Davis', 'sebastian.davis@example.com', 'tokendef');
INSERT INTO users (name, email, token) VALUES ('Charlotte Taylor', 'charlotte.taylor@example.com', 'tokenxyz');
