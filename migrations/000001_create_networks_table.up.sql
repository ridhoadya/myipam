CREATE TABLE networks (
  subnet_id SERIAL PRIMARY KEY,
  parent_id INT,
  name VARCHAR(255) NOT NULL,
  subnet VARCHAR(255) NOT NULL,
  CONSTRAINT fk_parent
    FOREIGN KEY(parent_id) REFERENCES networks(subnet_id)
);