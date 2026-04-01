-- ============================================================
-- F&F_DB — Mechanic Workshop Management Platform
-- Full Database Migration Script
-- Version: 1.0
-- ============================================================

-- Step 1: Create the database
-- Run this line separately if you haven't created the DB yet:
-- CREATE DATABASE "F&F_DB";

-- Step 2: Connect to the database before running the rest
-- \c "F&F_DB"

-- ============================================================
-- Step 3: Enable useful extensions
-- ============================================================
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";  -- For UUID primary keys

-- ============================================================
-- Step 4: Create ENUM types
-- ============================================================

-- User roles
CREATE TYPE user_role AS ENUM ('admin', 'mechanic');

-- Service status
CREATE TYPE service_status AS ENUM ('pending', 'in_progress', 'completed', 'cancelled');

-- Payment methods
CREATE TYPE payment_method AS ENUM ('cash', 'card', 'transfer', 'other');

-- ============================================================
-- Step 5: Create Table — users
-- ============================================================
CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name          VARCHAR(100)        NOT NULL,
    email         VARCHAR(150)        NOT NULL UNIQUE,
    password_hash VARCHAR(255)        NOT NULL,
    role          user_role           NOT NULL DEFAULT 'mechanic',
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at    TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- ============================================================
-- Step 6: Create Table — vehicles
-- ============================================================
CREATE TABLE vehicles (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    client_name   VARCHAR(100)        NOT NULL,
    phone         VARCHAR(20),
    vehicle_model VARCHAR(100)        NOT NULL,
    plate_number  VARCHAR(20)         NOT NULL UNIQUE,
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at    TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- ============================================================
-- Step 7: Create Table — services
-- ============================================================
CREATE TABLE services (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    vehicle_id  UUID                NOT NULL REFERENCES vehicles(id) ON DELETE CASCADE,
    description TEXT                NOT NULL,
    cost        NUMERIC(10, 2)      NOT NULL DEFAULT 0.00,
    status      service_status      NOT NULL DEFAULT 'pending',
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- ============================================================
-- Step 8: Create Table — payments
-- ============================================================
CREATE TABLE payments (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    service_id     UUID                NOT NULL REFERENCES services(id) ON DELETE CASCADE,
    amount         NUMERIC(10, 2)      NOT NULL,
    payment_method payment_method      NOT NULL DEFAULT 'cash',
    payment_date   DATE                NOT NULL DEFAULT CURRENT_DATE,
    notes          TEXT,
    created_at     TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- ============================================================
-- Step 9: Create Indexes for performance
-- ============================================================

-- Users
CREATE INDEX idx_users_email ON users(email);

-- Vehicles
CREATE INDEX idx_vehicles_plate_number  ON vehicles(plate_number);
CREATE INDEX idx_vehicles_client_name   ON vehicles(client_name);

-- Services
CREATE INDEX idx_services_vehicle_id ON services(vehicle_id);
CREATE INDEX idx_services_status     ON services(status);

-- Payments
CREATE INDEX idx_payments_service_id   ON payments(service_id);
CREATE INDEX idx_payments_payment_date ON payments(payment_date);

-- ============================================================
-- Step 10: Auto-update updated_at with a trigger function
-- ============================================================
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Attach trigger to users
CREATE TRIGGER trg_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Attach trigger to vehicles
CREATE TRIGGER trg_vehicles_updated_at
    BEFORE UPDATE ON vehicles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Attach trigger to services
CREATE TRIGGER trg_services_updated_at
    BEFORE UPDATE ON services
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ============================================================
-- Step 11: Seed data — default admin user (for testing)
-- Password: admin123 (bcrypt hash — replace in production!)
-- ============================================================
INSERT INTO users (name, email, password_hash, role)
VALUES (
    'Admin',
    'admin@ffworkshop.com',
    '$2a$12$placeholderHashReplaceThisBeforeGoingToProduction',
    'admin'
);

-- ============================================================
-- Done! F&F_DB is ready.
-- Tables: users, vehicles, services, payments
-- ============================================================
