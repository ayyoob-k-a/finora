-- ================================
-- Finora Database Setup Script
-- ================================
-- Run this as PostgreSQL superuser (postgres)

-- Create database
CREATE DATABASE finora_db;

-- Create user with password  
CREATE USER finora_user WITH ENCRYPTED PASSWORD 'finora_password';

-- Grant database privileges
GRANT ALL PRIVILEGES ON DATABASE finora_db TO finora_user;
GRANT USAGE ON SCHEMA public TO finora_user;
GRANT CREATE ON SCHEMA public TO finora_user;

-- Connect to finora_db database
\c finora_db

-- Grant table and sequence privileges
GRANT ALL ON ALL TABLES IN SCHEMA public TO finora_user;
GRANT ALL ON ALL SEQUENCES IN SCHEMA public TO finora_user;
GRANT ALL ON ALL FUNCTIONS IN SCHEMA public TO finora_user;

-- Set default privileges for future objects
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO finora_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO finora_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON FUNCTIONS TO finora_user;

-- Verify setup
SELECT 'Database setup completed successfully!' as result;
