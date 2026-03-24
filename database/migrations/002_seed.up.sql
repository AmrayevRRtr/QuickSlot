USE QuickSlot;

INSERT INTO users (email, password, role)
VALUES
    ('admin@test.com', '$2a$10$hashhashhash', 'ADMIN'),
    ('user@test.com', '$2a$10$hashhashhash', 'USER');

INSERT INTO organizations (name, owner_id)
VALUES ('Test Clinic', 1);

INSERT INTO employees (name, organization_id)
VALUES
    ('Dr. John', 1),
    ('Dr. Smith', 1);