-- Database initialization script for go-starter
-- This script runs when the PostgreSQL container starts for the first time

-- Create extensions if needed
-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create application tables (examples for future use)
-- CREATE TABLE IF NOT EXISTS users (
--     id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
--     email VARCHAR(255) UNIQUE NOT NULL,
--     password_hash VARCHAR(255) NOT NULL,
--     created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
--     updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
-- );

-- Create indexes for performance
-- CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Create application user (optional)
-- CREATE USER IF NOT EXISTS app_user WITH PASSWORD 'app_password';
-- GRANT CONNECT ON DATABASE go_starter TO app_user;
-- GRANT USAGE ON SCHEMA public TO app_user;
-- GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO app_user;
-- ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO app_user;

-- Insert initial data (optional)
-- INSERT INTO users (email, password_hash) VALUES 
-- ('admin@example.com', '$2a$10$placeholder_hash')
-- ON CONFLICT (email) DO NOTHING;

-- Create or update function to automatically update updated_at timestamp
-- CREATE OR REPLACE FUNCTION update_updated_at_column()
-- RETURNS TRIGGER AS $$
-- BEGIN
--     NEW.updated_at = NOW();
--     RETURN NEW;
-- END;
-- $$ language 'plpgsql';

-- Create triggers for updated_at (example)
-- CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
--     FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Log initialization
DO $$
BEGIN
    RAISE NOTICE 'Database initialized successfully for go-starter';
END $$;
