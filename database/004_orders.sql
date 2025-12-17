-- Migration: 004_orders.sql
-- Description: Complete orders system with payments, delivery, disputes
-- Date: 2024-12-17

-- =====================================================
-- ORDER STATUS ENUM
-- =====================================================
CREATE TYPE order_status AS ENUM (
    'PENDING',           -- Ordine creato, in attesa di pagamento
    'PAID',              -- Pagato, in attesa di consegna
    'PROCESSING',        -- Seller sta preparando
    'SHIPPED',           -- Spedito (per spedizione)
    'READY_FOR_PICKUP',  -- Pronto per ritiro
    'IN_TRANSIT',        -- In transito
    'DELIVERED',         -- Consegnato (tracking o QR)
    'COMPLETED',         -- Completato (dopo buffer 48h)
    'CANCELLED',         -- Annullato
    'REFUNDED',          -- Rimborsato
    'DISPUTED'           -- In contestazione
);

-- =====================================================
-- DELIVERY TYPE ENUM
-- =====================================================
CREATE TYPE delivery_type AS ENUM (
    'PICKUP',            -- Ritiro dal seller
    'SELLER_SHIPS',      -- Seller spedisce
    'BUYER_ARRANGES'     -- Buyer organizza ritiro
);

-- =====================================================
-- DISPUTE REASON ENUM
-- =====================================================
CREATE TYPE dispute_reason AS ENUM (
    'ITEM_NOT_RECEIVED',
    'ITEM_DAMAGED',
    'ITEM_NOT_AS_DESCRIBED',
    'SELLER_NO_SHOW',
    'BUYER_NO_SHOW',
    'SCAM_ATTEMPT',
    'OTHER'
);

-- =====================================================
-- DISPUTE STATUS ENUM
-- =====================================================
CREATE TYPE dispute_status AS ENUM (
    'OPEN',
    'SELLER_RESPONSE',
    'BUYER_REVIEW',
    'ADMIN_REVIEW',
    'RESOLVED_REFUND_FULL',
    'RESOLVED_REFUND_PARTIAL',
    'RESOLVED_PAYOUT_SELLER',
    'RESOLVED_SPLIT',
    'CLOSED'
);

-- =====================================================
-- ORDERS TABLE
-- =====================================================
CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- References
    buyer_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    seller_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,

    -- Order info
    quantity INT NOT NULL DEFAULT 1,
    unit_price DECIMAL(10,2) NOT NULL,
    shipping_cost DECIMAL(10,2) DEFAULT 0,
    total_amount DECIMAL(10,2) NOT NULL,

    -- Fees (calculated)
    platform_fee DECIMAL(10,2) DEFAULT 0,      -- 10% GecoGreen
    stripe_fee DECIMAL(10,2) DEFAULT 0,        -- ~1.4% + 0.25€
    seller_payout DECIMAL(10,2) DEFAULT 0,     -- What seller receives

    -- Status
    status order_status DEFAULT 'PENDING',
    delivery_type delivery_type NOT NULL,

    -- Delivery details (revealed after payment)
    pickup_location_id UUID REFERENCES seller_locations(id),
    pickup_address TEXT,
    pickup_instructions TEXT,
    pickup_deadline TIMESTAMP,

    -- Shipping details
    shipping_address TEXT,
    shipping_city VARCHAR(100),
    shipping_province VARCHAR(50),
    shipping_postal_code VARCHAR(10),
    shipping_country VARCHAR(50) DEFAULT 'IT',
    tracking_number VARCHAR(100),
    tracking_url VARCHAR(500),
    shipping_carrier VARCHAR(50),
    shipped_at TIMESTAMP,

    -- QR Code for pickup
    qr_code_token VARCHAR(100) UNIQUE,
    qr_code_expires_at TIMESTAMP,
    qr_scanned_at TIMESTAMP,

    -- Payment (Stripe)
    stripe_payment_intent_id VARCHAR(100),
    stripe_checkout_session_id VARCHAR(100),
    stripe_transfer_id VARCHAR(100),  -- For payout to seller
    paid_at TIMESTAMP,

    -- Payout
    payout_scheduled_at TIMESTAMP,
    payout_completed_at TIMESTAMP,
    payout_hold_reason TEXT,

    -- Impact tracking
    co2_saved DECIMAL(10,2) DEFAULT 0,
    water_saved DECIMAL(10,2) DEFAULT 0,
    eco_credits_buyer INT DEFAULT 0,
    eco_credits_seller INT DEFAULT 0,

    -- Notes
    buyer_notes TEXT,
    seller_notes TEXT,
    internal_notes TEXT,

    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP,
    cancelled_at TIMESTAMP,
    cancelled_by UUID REFERENCES users(id),
    cancellation_reason TEXT
);

-- Indexes
CREATE INDEX idx_orders_buyer ON orders(buyer_id, status);
CREATE INDEX idx_orders_seller ON orders(seller_id, status);
CREATE INDEX idx_orders_product ON orders(product_id);
CREATE INDEX idx_orders_status ON orders(status, created_at);
CREATE INDEX idx_orders_qr ON orders(qr_code_token) WHERE qr_code_token IS NOT NULL;
CREATE INDEX idx_orders_stripe ON orders(stripe_payment_intent_id);
CREATE INDEX idx_orders_payout ON orders(status, payout_scheduled_at)
    WHERE status = 'DELIVERED' AND payout_completed_at IS NULL;

-- =====================================================
-- ORDER STATUS HISTORY
-- =====================================================
CREATE TABLE IF NOT EXISTS order_status_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    old_status order_status,
    new_status order_status NOT NULL,
    changed_by UUID REFERENCES users(id),
    change_reason TEXT,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_order_history ON order_status_history(order_id, created_at);

-- =====================================================
-- DISPUTES TABLE
-- =====================================================
CREATE TABLE IF NOT EXISTS disputes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    opened_by UUID NOT NULL REFERENCES users(id),

    -- Dispute details
    reason dispute_reason NOT NULL,
    description TEXT NOT NULL,
    evidence_urls TEXT[], -- Array of image URLs

    -- Status
    status dispute_status DEFAULT 'OPEN',

    -- Response
    seller_response TEXT,
    seller_response_at TIMESTAMP,
    seller_evidence_urls TEXT[],

    -- Resolution
    resolved_by UUID REFERENCES users(id),
    resolution_notes TEXT,
    refund_amount DECIMAL(10,2),
    seller_payout_amount DECIMAL(10,2),

    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    resolved_at TIMESTAMP,

    -- Deadline for responses
    seller_response_deadline TIMESTAMP,
    admin_review_deadline TIMESTAMP
);

CREATE INDEX idx_disputes_order ON disputes(order_id);
CREATE INDEX idx_disputes_status ON disputes(status);
CREATE INDEX idx_disputes_deadline ON disputes(seller_response_deadline)
    WHERE status = 'OPEN';

-- =====================================================
-- USER STRIKES (for no-show penalties)
-- =====================================================
CREATE TABLE IF NOT EXISTS user_strikes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    order_id UUID REFERENCES orders(id) ON DELETE SET NULL,

    strike_type VARCHAR(50) NOT NULL, -- 'NO_SHOW', 'DISPUTE_LOST', 'FRAUD', etc.
    description TEXT,

    -- Strike can expire
    expires_at TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_strikes_user ON user_strikes(user_id, is_active);

-- =====================================================
-- CART TABLE (optional, for multi-item orders later)
-- =====================================================
CREATE TABLE IF NOT EXISTS cart_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    quantity INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(user_id, product_id)
);

CREATE INDEX idx_cart_user ON cart_items(user_id);

-- =====================================================
-- REVIEWS/FEEDBACK TABLE
-- =====================================================
CREATE TABLE IF NOT EXISTS order_reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,

    -- Who is reviewing whom
    reviewer_id UUID NOT NULL REFERENCES users(id),
    reviewed_id UUID NOT NULL REFERENCES users(id),

    -- Rating
    rating INT NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment TEXT,
    is_anonymous BOOLEAN DEFAULT FALSE,

    -- Moderation
    is_approved BOOLEAN DEFAULT TRUE,
    moderation_notes TEXT,

    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    -- One review per order per direction
    UNIQUE(order_id, reviewer_id)
);

CREATE INDEX idx_reviews_reviewed ON order_reviews(reviewed_id, is_approved);
CREATE INDEX idx_reviews_order ON order_reviews(order_id);

-- =====================================================
-- FUNCTIONS
-- =====================================================

-- Function to generate QR code token
CREATE OR REPLACE FUNCTION generate_qr_token()
RETURNS VARCHAR(100) AS $$
DECLARE
    token VARCHAR(100);
BEGIN
    token := encode(gen_random_bytes(32), 'hex');
    RETURN token;
END;
$$ LANGUAGE plpgsql;

-- Function to calculate fees
CREATE OR REPLACE FUNCTION calculate_order_fees(
    p_total_amount DECIMAL(10,2)
)
RETURNS TABLE (
    platform_fee DECIMAL(10,2),
    stripe_fee DECIMAL(10,2),
    seller_payout DECIMAL(10,2)
) AS $$
BEGIN
    -- Platform fee: 10%
    platform_fee := ROUND(p_total_amount * 0.10, 2);

    -- Stripe fee: 1.4% + 0.25€ (EU cards)
    stripe_fee := ROUND(p_total_amount * 0.014 + 0.25, 2);

    -- Seller gets the rest
    seller_payout := p_total_amount - platform_fee - stripe_fee;

    RETURN NEXT;
END;
$$ LANGUAGE plpgsql;

-- Trigger to update order timestamp
CREATE OR REPLACE FUNCTION update_order_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_order_timestamp
    BEFORE UPDATE ON orders
    FOR EACH ROW
    EXECUTE FUNCTION update_order_timestamp();

-- Trigger to log status changes
CREATE OR REPLACE FUNCTION log_order_status_change()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.status IS DISTINCT FROM NEW.status THEN
        INSERT INTO order_status_history (order_id, old_status, new_status)
        VALUES (NEW.id, OLD.status, NEW.status);
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_log_order_status
    AFTER UPDATE ON orders
    FOR EACH ROW
    EXECUTE FUNCTION log_order_status_change();

-- Trigger to update user rating after review
CREATE OR REPLACE FUNCTION update_user_rating()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE users
    SET
        rating_avg = (
            SELECT COALESCE(AVG(rating), 0)
            FROM order_reviews
            WHERE reviewed_id = NEW.reviewed_id AND is_approved = TRUE
        ),
        rating_count = (
            SELECT COUNT(*)
            FROM order_reviews
            WHERE reviewed_id = NEW.reviewed_id AND is_approved = TRUE
        ),
        updated_at = CURRENT_TIMESTAMP
    WHERE id = NEW.reviewed_id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_user_rating
    AFTER INSERT OR UPDATE ON order_reviews
    FOR EACH ROW
    EXECUTE FUNCTION update_user_rating();

-- Function to check if user can place order (strike check)
CREATE OR REPLACE FUNCTION can_user_order(p_user_id UUID)
RETURNS BOOLEAN AS $$
DECLARE
    active_strikes INT;
BEGIN
    SELECT COUNT(*) INTO active_strikes
    FROM user_strikes
    WHERE user_id = p_user_id
      AND is_active = TRUE
      AND (expires_at IS NULL OR expires_at > CURRENT_TIMESTAMP);

    -- Ban at 5+ strikes
    IF active_strikes >= 5 THEN
        RETURN FALSE;
    END IF;

    RETURN TRUE;
END;
$$ LANGUAGE plpgsql;

-- =====================================================
-- COMMENTS
-- =====================================================
COMMENT ON TABLE orders IS 'Main orders table with full payment and delivery tracking';
COMMENT ON TABLE disputes IS 'Dispute resolution system for problematic orders';
COMMENT ON TABLE user_strikes IS 'Strike system for no-shows and bad behavior';
COMMENT ON TABLE order_reviews IS 'Buyer/seller reviews after completed orders';
COMMENT ON TABLE cart_items IS 'Shopping cart for future multi-item orders';
