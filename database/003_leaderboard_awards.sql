-- Migration: 003_leaderboard_awards.sql
-- Description: Add tables for leaderboard, awards, and admin content tasks
-- Date: 2024-12-15

-- =====================================================
-- LEADERBOARD SNAPSHOTS
-- Stores periodic snapshots of user rankings
-- =====================================================
CREATE TABLE IF NOT EXISTS leaderboard_snapshots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,

    -- Period info
    period_type VARCHAR(20) NOT NULL CHECK (period_type IN ('WEEKLY', 'MONTHLY', 'YEARLY', 'ALLTIME')),
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,

    -- Metrics for this period
    total_co2_saved DECIMAL(10,2) DEFAULT 0,
    total_water_saved DECIMAL(10,2) DEFAULT 0,
    total_products_sold INT DEFAULT 0,
    total_orders INT DEFAULT 0,
    total_revenue DECIMAL(10,2) DEFAULT 0,

    -- Ranking
    rank INT,

    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_leaderboard_period ON leaderboard_snapshots(period_type, period_start);
CREATE INDEX idx_leaderboard_user ON leaderboard_snapshots(user_id);
CREATE INDEX idx_leaderboard_rank ON leaderboard_snapshots(period_type, period_start, rank);

-- =====================================================
-- AWARDS
-- Stores awards given to users (Eco-Champion, New Entry, etc.)
-- =====================================================
CREATE TYPE award_type AS ENUM (
    'ECO_CHAMPION',      -- 1st place monthly
    'ECO_RUNNER_UP',     -- 2nd place monthly
    'ECO_THIRD',         -- 3rd place monthly
    'NEW_ENTRY',         -- Best new seller of the month
    'RECORD_BREAKER',    -- Broke a record
    'TOP_WEEK',          -- Weekly top seller
    'ECO_LEGEND',        -- All-time achievement
    'MILESTONE'          -- Reached a milestone (10, 50, 100 products)
);

CREATE TYPE interview_status AS ENUM (
    'NOT_REQUIRED',
    'PENDING',
    'CONTACTED',
    'SCHEDULED',
    'RECORDED',
    'PUBLISHED',
    'DECLINED'
);

CREATE TABLE IF NOT EXISTS awards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,

    -- Award info
    award_type award_type NOT NULL,
    period_type VARCHAR(20), -- 'WEEKLY', 'MONTHLY', 'YEARLY'
    period_start DATE,
    period_end DATE,

    -- Details
    title VARCHAR(200) NOT NULL,
    description TEXT,
    badge_url VARCHAR(500),

    -- Stats at time of award
    co2_saved DECIMAL(10,2),
    products_count INT,

    -- Content/Interview
    youtube_url VARCHAR(500),
    interview_status interview_status DEFAULT 'PENDING',
    interview_scheduled_at TIMESTAMP,
    interview_notes TEXT,

    -- Visibility
    is_featured BOOLEAN DEFAULT FALSE,
    is_public BOOLEAN DEFAULT TRUE,

    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    published_at TIMESTAMP
);

CREATE INDEX idx_awards_user ON awards(user_id);
CREATE INDEX idx_awards_type ON awards(award_type, period_start);
CREATE INDEX idx_awards_featured ON awards(is_featured) WHERE is_featured = TRUE;
CREATE INDEX idx_awards_interview ON awards(interview_status) WHERE interview_status IN ('PENDING', 'CONTACTED', 'SCHEDULED');

-- =====================================================
-- ADMIN CONTENT TASKS
-- Reminder system for admin to create content
-- =====================================================
CREATE TYPE task_status AS ENUM (
    'PENDING',
    'IN_PROGRESS',
    'COMPLETED',
    'SKIPPED',
    'CANCELLED'
);

CREATE TYPE task_priority AS ENUM (
    'URGENT',
    'HIGH',
    'NORMAL',
    'LOW'
);

CREATE TYPE task_type AS ENUM (
    'INTERVIEW_CONTACT',     -- Contact winner for interview
    'INTERVIEW_SCHEDULE',    -- Schedule the interview
    'INTERVIEW_RECORD',      -- Record the interview
    'YOUTUBE_EDIT',          -- Edit video
    'YOUTUBE_PUBLISH',       -- Publish to YouTube
    'SOCIAL_POST',           -- Post on social media
    'HALL_OF_FAME_UPDATE',   -- Update Hall of Fame page
    'EMAIL_WINNER',          -- Send email to winner
    'BADGE_ASSIGN',          -- Assign badge
    'OTHER'
);

CREATE TABLE IF NOT EXISTS admin_content_tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- References
    award_id UUID REFERENCES awards(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL, -- The user this task is about

    -- Task info
    task_type task_type NOT NULL,
    title VARCHAR(200) NOT NULL,
    description TEXT,

    -- Status
    status task_status DEFAULT 'PENDING',
    priority task_priority DEFAULT 'NORMAL',
    due_date DATE,

    -- Completion
    completed_at TIMESTAMP,
    completed_by UUID REFERENCES users(id),
    notes TEXT,

    -- Auto-generated
    is_auto_generated BOOLEAN DEFAULT TRUE,

    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_admin_tasks_status ON admin_content_tasks(status, due_date);
CREATE INDEX idx_admin_tasks_award ON admin_content_tasks(award_id);
CREATE INDEX idx_admin_tasks_pending ON admin_content_tasks(status, priority, due_date)
    WHERE status IN ('PENDING', 'IN_PROGRESS');

-- =====================================================
-- IMPACT LOGS
-- Detailed log of all impact/credits actions
-- =====================================================
CREATE TYPE impact_action_type AS ENUM (
    'PURCHASE',
    'SALE',
    'GIFT_GIVEN',
    'GIFT_RECEIVED',
    'FIRST_PURCHASE',
    'FIRST_SALE',
    'PICKUP_BONUS',
    'LAST_CHANCE_BONUS',
    'REVIEW_BONUS',
    'REFERRAL_BONUS',
    'MILESTONE_BONUS',
    'SOCIAL_SHARE',
    'PROFILE_COMPLETE',
    'REDEEM_BOOST',
    'REDEEM_TREE',
    'REDEEM_BADGE',
    'POINTS_EXPIRED',
    'ADMIN_ADJUSTMENT'
);

CREATE TABLE IF NOT EXISTS impact_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    order_id UUID REFERENCES orders(id) ON DELETE SET NULL,

    -- Action
    action_type impact_action_type NOT NULL,

    -- Impact metrics
    co2_saved DECIMAL(10,2) DEFAULT 0,
    water_saved DECIMAL(10,2) DEFAULT 0,

    -- Credits
    eco_credits_earned INT DEFAULT 0,
    eco_credits_spent INT DEFAULT 0,
    eco_credits_balance INT, -- Balance after this transaction

    -- Description
    description TEXT,
    metadata JSONB, -- Additional data (e.g., referral user_id, tree certificate, etc.)

    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_impact_logs_user ON impact_logs(user_id);
CREATE INDEX idx_impact_logs_order ON impact_logs(order_id);
CREATE INDEX idx_impact_logs_date ON impact_logs(created_at);
CREATE INDEX idx_impact_logs_type ON impact_logs(action_type, created_at);

-- =====================================================
-- TREES PLANTED
-- Track trees planted via Tree-Nation
-- =====================================================
CREATE TABLE IF NOT EXISTS trees_planted (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Who planted
    user_id UUID REFERENCES users(id) ON DELETE SET NULL, -- NULL if platform planted
    impact_log_id UUID REFERENCES impact_logs(id) ON DELETE SET NULL,

    -- Tree-Nation info
    tree_nation_id VARCHAR(100), -- ID from Tree-Nation API
    tree_nation_url VARCHAR(500), -- Certificate URL
    species VARCHAR(100),
    location VARCHAR(200), -- Where planted (country/region)

    -- If platform planted (every 100 orders)
    is_platform_tree BOOLEAN DEFAULT FALSE,
    trigger_order_count INT, -- e.g., 100, 200, 300...

    -- Timestamps
    planted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_trees_user ON trees_planted(user_id);
CREATE INDEX idx_trees_platform ON trees_planted(is_platform_tree);

-- =====================================================
-- CSR REPORTS
-- Generated sustainability reports for businesses
-- =====================================================
CREATE TABLE IF NOT EXISTS csr_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    seller_id UUID REFERENCES users(id) ON DELETE CASCADE,

    -- Period
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    report_type VARCHAR(20) DEFAULT 'YEARLY', -- 'MONTHLY', 'QUARTERLY', 'YEARLY'

    -- Aggregated stats
    total_items INT DEFAULT 0,
    total_co2_kg DECIMAL(10,2) DEFAULT 0,
    total_water_l DECIMAL(10,2) DEFAULT 0,
    total_waste_kg DECIMAL(10,2) DEFAULT 0,
    total_value DECIMAL(10,2) DEFAULT 0,

    -- Category breakdown (JSON)
    category_breakdown JSONB,

    -- Generated files
    pdf_url VARCHAR(500),
    report_code VARCHAR(50) UNIQUE,

    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_csr_seller ON csr_reports(seller_id);
CREATE INDEX idx_csr_period ON csr_reports(period_start, period_end);

-- =====================================================
-- HELPER FUNCTIONS
-- =====================================================

-- Function to get current leaderboard
CREATE OR REPLACE FUNCTION get_current_leaderboard(
    p_period_type VARCHAR(20),
    p_limit INT DEFAULT 10
)
RETURNS TABLE (
    rank BIGINT,
    user_id UUID,
    business_name VARCHAR,
    first_name VARCHAR,
    last_name VARCHAR,
    account_type VARCHAR,
    city VARCHAR,
    total_co2_saved DECIMAL,
    total_products_sold BIGINT,
    avatar_url VARCHAR
) AS $$
BEGIN
    RETURN QUERY
    WITH period_dates AS (
        SELECT
            CASE p_period_type
                WHEN 'WEEKLY' THEN date_trunc('week', CURRENT_DATE)::DATE
                WHEN 'MONTHLY' THEN date_trunc('month', CURRENT_DATE)::DATE
                WHEN 'YEARLY' THEN date_trunc('year', CURRENT_DATE)::DATE
                ELSE '2024-01-01'::DATE -- ALLTIME
            END as start_date,
            CURRENT_DATE as end_date
    ),
    seller_stats AS (
        SELECT
            p.seller_id,
            COALESCE(SUM(
                CASE
                    WHEN c.estimated_co2_kg IS NOT NULL THEN c.estimated_co2_kg * p.quantity_available
                    ELSE 2.0 * p.quantity_available -- Default 2kg if no category
                END
            ), 0) as co2_saved,
            COUNT(DISTINCT o.id) as products_sold
        FROM products p
        LEFT JOIN categories c ON p.category_id = c.id
        LEFT JOIN orders o ON o.product_id = p.id AND o.status = 'COMPLETED'
        CROSS JOIN period_dates pd
        WHERE o.created_at >= pd.start_date
        GROUP BY p.seller_id
    )
    SELECT
        ROW_NUMBER() OVER (ORDER BY ss.co2_saved DESC) as rank,
        u.id as user_id,
        u.business_name,
        u.first_name,
        u.last_name,
        u.account_type::VARCHAR,
        u.city,
        ss.co2_saved as total_co2_saved,
        ss.products_sold as total_products_sold,
        u.avatar_url
    FROM seller_stats ss
    JOIN users u ON u.id = ss.seller_id
    WHERE u.status = 'ACTIVE'
    ORDER BY ss.co2_saved DESC
    LIMIT p_limit;
END;
$$ LANGUAGE plpgsql;

-- Function to create monthly awards automatically
CREATE OR REPLACE FUNCTION create_monthly_awards()
RETURNS void AS $$
DECLARE
    v_period_start DATE;
    v_period_end DATE;
    v_user RECORD;
    v_rank INT := 0;
BEGIN
    -- Calculate last month
    v_period_start := date_trunc('month', CURRENT_DATE - INTERVAL '1 month')::DATE;
    v_period_end := (date_trunc('month', CURRENT_DATE) - INTERVAL '1 day')::DATE;

    -- Get top 3 sellers
    FOR v_user IN
        SELECT * FROM get_current_leaderboard('MONTHLY', 3)
    LOOP
        v_rank := v_rank + 1;

        -- Create award
        INSERT INTO awards (
            user_id, award_type, period_type, period_start, period_end,
            title, description, co2_saved, products_count, interview_status
        ) VALUES (
            v_user.user_id,
            CASE v_rank
                WHEN 1 THEN 'ECO_CHAMPION'::award_type
                WHEN 2 THEN 'ECO_RUNNER_UP'::award_type
                WHEN 3 THEN 'ECO_THIRD'::award_type
            END,
            'MONTHLY',
            v_period_start,
            v_period_end,
            CASE v_rank
                WHEN 1 THEN 'Eco-Champion del mese'
                WHEN 2 THEN 'Secondo classificato del mese'
                WHEN 3 THEN 'Terzo classificato del mese'
            END,
            COALESCE(v_user.business_name, v_user.first_name || ' ' || v_user.last_name) ||
            ' ha risparmiato ' || v_user.total_co2_saved || ' kg di CO₂',
            v_user.total_co2_saved,
            v_user.total_products_sold,
            CASE WHEN v_rank <= 2 THEN 'PENDING'::interview_status ELSE 'NOT_REQUIRED'::interview_status END
        );
    END LOOP;
END;
$$ LANGUAGE plpgsql;

-- Function to auto-generate admin tasks for new awards
CREATE OR REPLACE FUNCTION auto_create_admin_tasks()
RETURNS TRIGGER AS $$
BEGIN
    -- Only create tasks for awards that need interviews
    IF NEW.interview_status = 'PENDING' THEN
        -- Task: Contact winner
        INSERT INTO admin_content_tasks (
            award_id, user_id, task_type, title, description, priority, due_date
        ) VALUES (
            NEW.id,
            NEW.user_id,
            'INTERVIEW_CONTACT',
            'Contattare ' || NEW.title,
            'Contattare il vincitore per organizzare intervista YouTube',
            'HIGH',
            CURRENT_DATE + INTERVAL '3 days'
        );

        -- Task: YouTube publish (due later)
        INSERT INTO admin_content_tasks (
            award_id, user_id, task_type, title, description, priority, due_date
        ) VALUES (
            NEW.id,
            NEW.user_id,
            'YOUTUBE_PUBLISH',
            'Pubblicare video ' || NEW.title,
            'Pubblicare intervista su YouTube e aggiornare Hall of Fame',
            'NORMAL',
            CURRENT_DATE + INTERVAL '10 days'
        );
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_auto_create_admin_tasks
    AFTER INSERT ON awards
    FOR EACH ROW
    EXECUTE FUNCTION auto_create_admin_tasks();

-- =====================================================
-- UPDATE users table (add eco_credits if not exists)
-- =====================================================
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns
                   WHERE table_name = 'users' AND column_name = 'eco_credits') THEN
        ALTER TABLE users ADD COLUMN eco_credits INT DEFAULT 0;
    END IF;

    IF NOT EXISTS (SELECT 1 FROM information_schema.columns
                   WHERE table_name = 'users' AND column_name = 'eco_level') THEN
        ALTER TABLE users ADD COLUMN eco_level VARCHAR(50) DEFAULT 'Germoglio';
    END IF;

    IF NOT EXISTS (SELECT 1 FROM information_schema.columns
                   WHERE table_name = 'users' AND column_name = 'total_co2_saved') THEN
        ALTER TABLE users ADD COLUMN total_co2_saved DECIMAL(10,2) DEFAULT 0;
    END IF;

    IF NOT EXISTS (SELECT 1 FROM information_schema.columns
                   WHERE table_name = 'users' AND column_name = 'total_water_saved') THEN
        ALTER TABLE users ADD COLUMN total_water_saved DECIMAL(10,2) DEFAULT 0;
    END IF;
END $$;

-- =====================================================
-- UPDATE categories table (add impact values if not exists)
-- =====================================================
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns
                   WHERE table_name = 'categories' AND column_name = 'estimated_co2_kg') THEN
        ALTER TABLE categories ADD COLUMN estimated_co2_kg DECIMAL(10,2) DEFAULT 2;
    END IF;

    IF NOT EXISTS (SELECT 1 FROM information_schema.columns
                   WHERE table_name = 'categories' AND column_name = 'estimated_water_l') THEN
        ALTER TABLE categories ADD COLUMN estimated_water_l DECIMAL(10,2) DEFAULT 500;
    END IF;

    IF NOT EXISTS (SELECT 1 FROM information_schema.columns
                   WHERE table_name = 'categories' AND column_name = 'estimated_waste_kg') THEN
        ALTER TABLE categories ADD COLUMN estimated_waste_kg DECIMAL(10,2) DEFAULT 1;
    END IF;
END $$;

-- =====================================================
-- SEED: Update category impact values
-- =====================================================
UPDATE categories SET estimated_co2_kg = 2, estimated_water_l = 500, estimated_waste_kg = 1
WHERE slug = 'alimentari-freschi' OR name ILIKE '%freschi%';

UPDATE categories SET estimated_co2_kg = 1, estimated_water_l = 200, estimated_waste_kg = 0.5
WHERE slug = 'alimentari-confezionati' OR name ILIKE '%confezionat%';

UPDATE categories SET estimated_co2_kg = 0.5, estimated_water_l = 100, estimated_waste_kg = 0.3
WHERE slug = 'bevande' OR name ILIKE '%bevand%';

UPDATE categories SET estimated_co2_kg = 15, estimated_water_l = 500, estimated_waste_kg = 5
WHERE slug = 'materiali-tecnici' OR name ILIKE '%tecnic%';

UPDATE categories SET estimated_co2_kg = 30, estimated_water_l = 2000, estimated_waste_kg = 10
WHERE slug = 'automotive' OR name ILIKE '%auto%';

COMMENT ON TABLE leaderboard_snapshots IS 'Stores periodic snapshots of user rankings for hall of fame';
COMMENT ON TABLE awards IS 'Awards given to users (Eco-Champion, New Entry, etc.)';
COMMENT ON TABLE admin_content_tasks IS 'Task management for admin to create YouTube content';
COMMENT ON TABLE impact_logs IS 'Detailed log of all eco-credits transactions';
COMMENT ON TABLE trees_planted IS 'Track trees planted via Tree-Nation integration';
COMMENT ON TABLE csr_reports IS 'Generated sustainability reports for businesses';
