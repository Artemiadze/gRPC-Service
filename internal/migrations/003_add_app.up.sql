INSERT INTO apps (id, name, secret)
VALUES (1, 'test', 'test-secret')
ON CONFLICT (name) DO NOTHING;
