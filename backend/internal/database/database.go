package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// DB holds database connections
type DB struct {
	Pool  *pgxpool.Pool
	Redis *redis.Client
}

// New creates a new database connection
func New(databaseURL, redisURL string) (*DB, error) {
	// PostgreSQL connection pool
	poolConfig, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database URL: %w", err)
	}

	// Connection pool settings
	poolConfig.MaxConns = 25
	poolConfig.MinConns = 5
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Test PostgreSQL connection
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	// Redis connection
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to parse redis URL: %w", err)
	}

	redisClient := redis.NewClient(opt)

	// Test Redis connection
	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to ping redis: %w", err)
	}

	return &DB{
		Pool:  pool,
		Redis: redisClient,
	}, nil
}

// Close closes all database connections
func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
	if db.Redis != nil {
		db.Redis.Close()
	}
}

// Health checks database health
func (db *DB) Health(ctx context.Context) map[string]string {
	health := make(map[string]string)

	// Check PostgreSQL
	if err := db.Pool.Ping(ctx); err != nil {
		health["postgres"] = fmt.Sprintf("error: %v", err)
	} else {
		health["postgres"] = "ok"
	}

	// Check Redis
	if _, err := db.Redis.Ping(ctx).Result(); err != nil {
		health["redis"] = fmt.Sprintf("error: %v", err)
	} else {
		health["redis"] = "ok"
	}

	return health
}

// Migrate creates database schema if not exists
func (db *DB) Migrate(ctx context.Context) error {
	log.Println("🔄 Running database migrations...")

	// Check if migration already done (users table exists)
	var exists bool
	err := db.Pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT FROM information_schema.tables
			WHERE table_schema = 'public' AND table_name = 'users'
		)
	`).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check migration status: %w", err)
	}

	if !exists {
		log.Println("📦 Creating initial database schema...")
		tx, err := db.Pool.Begin(ctx)
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}
		defer tx.Rollback(ctx)

		_, err = tx.Exec(ctx, migrationSQL)
		if err != nil {
			return fmt.Errorf("failed to execute migration: %w", err)
		}

		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("failed to commit migration: %w", err)
		}
		log.Println("✅ Initial schema created")
	}

	// Run incremental migrations
	if err := db.runIncrementalMigrations(ctx); err != nil {
		return fmt.Errorf("failed to run incremental migrations: %w", err)
	}

	log.Println("✅ Database migration completed successfully")
	return nil
}

// runIncrementalMigrations runs migrations that add new features
func (db *DB) runIncrementalMigrations(ctx context.Context) error {
	// Migration 1: Add account_type and business fields
	var hasAccountType bool
	err := db.Pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT FROM information_schema.columns
			WHERE table_name = 'users' AND column_name = 'account_type'
		)
	`).Scan(&hasAccountType)
	if err != nil {
		return err
	}

	if !hasAccountType {
		log.Println("📦 Adding account_type and business fields...")
		_, err = db.Pool.Exec(ctx, migrationV2SQL)
		if err != nil {
			return fmt.Errorf("migration v2 failed: %w", err)
		}
		log.Println("✅ Migration v2 completed")
	}

	return nil
}

// Migration V2: Account type, business fields, social links, image moderation
const migrationV2SQL = `
-- Account type enum
DO $$ BEGIN
    CREATE TYPE account_type AS ENUM ('PRIVATE', 'BUSINESS');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Image moderation status
DO $$ BEGIN
    CREATE TYPE moderation_status AS ENUM ('PENDING', 'APPROVED', 'REJECTED');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Add new columns to users
ALTER TABLE users ADD COLUMN IF NOT EXISTS account_type account_type DEFAULT 'PRIVATE';
ALTER TABLE users ADD COLUMN IF NOT EXISTS business_name VARCHAR(255);
ALTER TABLE users ADD COLUMN IF NOT EXISTS vat_number VARCHAR(20);
ALTER TABLE users ADD COLUMN IF NOT EXISTS has_multiple_locations BOOLEAN DEFAULT FALSE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS social_links JSONB DEFAULT '{}';
ALTER TABLE users ADD COLUMN IF NOT EXISTS business_photos JSONB DEFAULT '[]';
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_admin BOOLEAN DEFAULT FALSE;

-- Update seller_locations to work for all users (rename seller_id to user_id conceptually but keep for compatibility)
-- Add user_id as alias
DO $$ BEGIN
    ALTER TABLE seller_locations RENAME COLUMN seller_id TO user_id;
EXCEPTION
    WHEN undefined_column THEN null;
END $$;

-- Image moderation queue
CREATE TABLE IF NOT EXISTS image_reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    product_id UUID REFERENCES products(id) ON DELETE CASCADE,
    image_url VARCHAR(500) NOT NULL,
    image_type VARCHAR(20) NOT NULL, -- 'PRODUCT', 'PROFILE', 'BUSINESS'

    -- OCR results
    detected_text TEXT,
    detected_phone BOOLEAN DEFAULT FALSE,
    detected_email BOOLEAN DEFAULT FALSE,
    detected_url BOOLEAN DEFAULT FALSE,
    confidence_score DECIMAL(5,4),

    -- Moderation
    status moderation_status DEFAULT 'PENDING',
    reviewed_by UUID REFERENCES users(id),
    reviewed_at TIMESTAMP,
    rejection_reason TEXT,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index for pending reviews
CREATE INDEX IF NOT EXISTS idx_image_reviews_status ON image_reviews(status);
CREATE INDEX IF NOT EXISTS idx_image_reviews_user ON image_reviews(user_id);

-- Index for business accounts
CREATE INDEX IF NOT EXISTS idx_users_account_type ON users(account_type);
CREATE INDEX IF NOT EXISTS idx_users_vat ON users(vat_number);
`

const migrationSQL = `
-- ============================================
-- ENUM TYPES
-- ============================================

CREATE TYPE user_role AS ENUM ('BUYER', 'SELLER', 'ADMIN');
CREATE TYPE user_status AS ENUM ('PENDING', 'ACTIVE', 'SUSPENDED', 'BANNED');
CREATE TYPE listing_type AS ENUM ('SALE', 'GIFT');
CREATE TYPE shipping_method AS ENUM ('PICKUP', 'SELLER_SHIPS', 'BUYER_ARRANGES', 'PLATFORM_MANAGED');
CREATE TYPE product_status AS ENUM ('DRAFT', 'ACTIVE', 'SOLD', 'EXPIRED', 'DELETED');
CREATE TYPE order_status AS ENUM ('CREATED', 'AWAITING_PAYMENT', 'PAID', 'READY_PICKUP', 'SHIPPED', 'DELIVERED', 'COMPLETED', 'DISPUTED', 'CANCELLED', 'REFUNDED');

-- ============================================
-- TABLES
-- ============================================

-- UTENTI
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone VARCHAR(20),
    city VARCHAR(100),
    province VARCHAR(50),
    postal_code VARCHAR(10),
    country VARCHAR(2) DEFAULT 'IT',
    roles user_role[] DEFAULT '{BUYER}',
    status user_status DEFAULT 'PENDING',
    email_verified BOOLEAN DEFAULT FALSE,
    email_verification_token VARCHAR(100),
    email_verification_expires TIMESTAMP,
    password_reset_token VARCHAR(100),
    password_reset_expires TIMESTAMP,
    google_id VARCHAR(100) UNIQUE,
    apple_id VARCHAR(100) UNIQUE,
    facebook_id VARCHAR(100) UNIQUE,
    oauth_provider VARCHAR(20),
    avatar_url VARCHAR(500),
    stripe_customer_id VARCHAR(50),
    stripe_account_id VARCHAR(50),
    stripe_onboarding_complete BOOLEAN DEFAULT FALSE,
    rating_avg DECIMAL(3,2) DEFAULT 0,
    rating_count INTEGER DEFAULT 0,
    strike_count INTEGER DEFAULT 0,
    total_co2_saved DECIMAL(10,2) DEFAULT 0,
    total_water_saved DECIMAL(10,2) DEFAULT 0,
    eco_credits INT DEFAULT 0,
    eco_level VARCHAR(50) DEFAULT 'Germoglio',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login_at TIMESTAMP,
    deleted_at TIMESTAMP
);

-- SESSIONI
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    refresh_token VARCHAR(500) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    device_info JSONB,
    ip_address INET,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_used_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- CATEGORIE
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    icon VARCHAR(50),
    parent_id UUID REFERENCES categories(id),
    sort_order INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    estimated_co2_kg DECIMAL(10,2) DEFAULT 0,
    estimated_water_l DECIMAL(10,2) DEFAULT 0,
    estimated_waste_kg DECIMAL(10,2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- SELLER LOCATIONS
CREATE TABLE seller_locations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    seller_id UUID REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    is_primary BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    address_street VARCHAR(255) NOT NULL,
    address_city VARCHAR(100) NOT NULL,
    address_province VARCHAR(50),
    address_postal_code VARCHAR(10) NOT NULL,
    address_country VARCHAR(2) DEFAULT 'IT',
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    phone VARCHAR(20),
    email VARCHAR(255),
    pickup_hours JSONB,
    pickup_instructions TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- PRODOTTI
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    seller_id UUID REFERENCES users(id) ON DELETE CASCADE,
    category_id UUID REFERENCES categories(id),
    location_id UUID REFERENCES seller_locations(id),
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    original_price DECIMAL(10, 2),
    listing_type listing_type DEFAULT 'SALE',
    shipping_method shipping_method DEFAULT 'PICKUP',
    shipping_cost DECIMAL(10, 2) DEFAULT 0,
    quantity INTEGER DEFAULT 1,
    quantity_available INTEGER DEFAULT 1,
    expiry_date DATE,
    expiry_photo_url VARCHAR(500),
    is_dutch_auction BOOLEAN DEFAULT FALSE,
    dutch_start_price DECIMAL(10, 2),
    dutch_decrease_amount DECIMAL(10, 2),
    dutch_decrease_hours INT DEFAULT 24,
    dutch_min_price DECIMAL(10, 2),
    dutch_started_at TIMESTAMP,
    weight_kg DECIMAL(10, 2),
    city VARCHAR(100),
    province VARCHAR(50),
    postal_code VARCHAR(10),
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    images JSONB DEFAULT '[]',
    status product_status DEFAULT 'DRAFT',
    slug VARCHAR(255),
    view_count INTEGER DEFAULT 0,
    favorite_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    published_at TIMESTAMP,
    deleted_at TIMESTAMP
);

-- PREFERITI
CREATE TABLE favorites (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    product_id UUID REFERENCES products(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, product_id)
);

-- ============================================
-- INDEXES
-- ============================================

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_google ON users(google_id);
CREATE INDEX idx_sessions_user ON sessions(user_id);
CREATE INDEX idx_sessions_token ON sessions(refresh_token);
CREATE INDEX idx_products_seller ON products(seller_id);
CREATE INDEX idx_products_category ON products(category_id);
CREATE INDEX idx_products_status ON products(status);
CREATE INDEX idx_products_city ON products(city);
CREATE INDEX idx_products_created ON products(created_at DESC);
CREATE INDEX idx_seller_locations_seller ON seller_locations(seller_id);

-- ============================================
-- TRIGGER
-- ============================================

CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at();
CREATE TRIGGER products_updated_at BEFORE UPDATE ON products FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- ============================================
-- SEED DATA: Categorie
-- ============================================

INSERT INTO categories (id, name, slug, description, icon, sort_order) VALUES
('10000000-0000-0000-0000-000000000001', 'Alimentari Freschi', 'alimentari-freschi', 'Latticini, carne, pesce, uova, salumi, gastronomia', 'fresh', 1),
('10000000-0000-0000-0000-000000000002', 'Alimentari Confezionati', 'alimentari-confezionati', 'Conserve, snack, biscotti, salse, cereali', 'packaged', 2),
('10000000-0000-0000-0000-000000000003', 'Bevande', 'bevande', 'Latte UHT, succhi, birra, bibite, bevande vegetali', 'drinks', 3),
('10000000-0000-0000-0000-000000000004', 'Frutta e Verdura', 'frutta-verdura', 'Fresca e quarta gamma', 'produce', 4),
('10000000-0000-0000-0000-000000000005', 'Surgelati', 'surgelati', 'Pesce, carne, verdure, pizze, gelati', 'frozen', 5),
('10000000-0000-0000-0000-000000000006', 'Cosmetici', 'cosmetici', 'Creme, trucco, solari, profumi', 'cosmetics', 6),
('10000000-0000-0000-0000-000000000007', 'Cura Persona', 'cura-persona', 'Shampoo, dentifrici, igiene', 'personal-care', 7),
('10000000-0000-0000-0000-000000000008', 'Detergenza Casa', 'detergenza', 'Detersivi, candeggina, anticalcare', 'cleaning', 8),
('10000000-0000-0000-0000-000000000009', 'Pet Food', 'pet-food', 'Cibo e snack per animali', 'pets', 9),
('10000000-0000-0000-0000-000000000010', 'Giardinaggio', 'giardinaggio', 'Fertilizzanti, semi, fitofarmaci', 'garden', 10),
('10000000-0000-0000-0000-000000000011', 'Materiali Tecnici', 'materiali-tecnici', 'Vernici, colle, siliconi, resine', 'technical', 11),
('10000000-0000-0000-0000-000000000012', 'Automotive', 'automotive', 'Pneumatici, olio motore, accessori', 'automotive', 12),
('10000000-0000-0000-0000-000000000013', 'Sicurezza e DPI', 'sicurezza-dpi', 'Estintori, caschi, imbragature', 'safety', 13),
('10000000-0000-0000-0000-000000000014', 'HORECA', 'horeca', 'Ristorazione, catering, monoporzioni', 'horeca', 14);
`
