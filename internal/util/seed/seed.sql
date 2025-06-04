-- Seed Users
INSERT INTO users (username, email, password, phone, first_name, last_name, is_admin, created_at, updated_at) 
VALUES 
    ('admin', 'admin@example.com', '1234', '0891234567', 'Admin', 'User', true, NOW(), NOW()),
    ('john_doe', 'john@example.com', 'johndoe', '0891234568', 'John', 'Doe', false, NOW(), NOW()),
    ('jane_smith', 'jane@example.com', 'janesmith', '0891234569', 'Jane', 'Smith', false, NOW(), NOW()),
    ('bob_wilson', 'bob@example.com', 'bobwilson', '0891234570', 'Bob', 'Wilson', false, NOW(), NOW()),
    ('alice_brown', 'alice@example.com', 'alicebrown', '0891234571', 'Alice', 'Brown', false, NOW(), NOW());

-- Seed Items
INSERT INTO items (name, description, available_amount, total_amount, status, created_at, updated_at)
VALUES
    ('Laptop Dell XPS 13', 'High-performance laptop for development', 5, 5, 'available', NOW(), NOW()),
    ('iPad Pro 12.9', 'Latest iPad Pro with M1 chip', 3, 3, 'lost', NOW(), NOW()),
    ('MacBook Pro M1', '16-inch MacBook Pro with M1 Pro chip', 2, 2, 'borrowed', NOW(), NOW()),
    ('ThinkPad X1 Carbon', 'Lenovo ThinkPad business laptop', 4, 4, 'available', NOW(), NOW()),
    ('Surface Pro 8', 'Microsoft Surface Pro tablet/laptop', 3, 3, 'maintenance', NOW(), NOW()),
    ('iPhone 13 Pro', 'Test device for iOS development', 2, 2, 'available', NOW(), NOW()),
    ('Android Pixel 6', 'Test device for Android development', 2, 2, 'borrowed', NOW(), NOW()),
    ('Monitor Dell U2719D', '27-inch 4K USB-C monitor', 6, 6, 'available', NOW(), NOW()),
    ('Logitech MX Master 3', 'Wireless premium mouse', 10, 10, 'available', NOW(), NOW()),
    ('Herman Miller Aeron', 'Ergonomic office chair', 8, 8, 'available', NOW(), NOW());