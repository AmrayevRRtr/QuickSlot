USE QuickSlot;

-- users
INSERT INTO users (email, password, role)
VALUES
    ('user@test.com', '$2a$10$abcdefghijklmnopqrstuv', 'USER');

-- organization
INSERT INTO organizations (name, owner_id)
VALUES ('Test Clinic', 1);

-- employees
INSERT INTO employees (name, organization_id)
VALUES
    ('Dr. John', 1);

INSERT INTO employees (name, organization_id)
VALUES
    ('Dr. Dre', 1);

-- slots
INSERT INTO time_slots (employee_id, start_time, end_time, is_booked)
VALUES
    (1, '2026-03-25 10:00:00', '2026-03-25 10:30:00', FALSE),
    ( 1, '2026-03-25 10:30:00', '2026-03-25 11:00:00', FALSE);