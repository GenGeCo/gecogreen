-- ============================================
-- GECOGREEN DATABASE SCHEMA
-- PostgreSQL 16+
-- ============================================
-- Eseguire in ordine: ENUM, TABLES, INDEXES
-- ============================================

-- ============================================
-- ENUM TYPES
-- ============================================

-- Ruoli utente
CREATE TYPE user_role AS ENUM (
    'BUYER',      -- Cliente che compra
    'SELLER',     -- Venditore verificato
    'ADMIN'       -- Amministratore
);

-- Stato utente
CREATE TYPE user_status AS ENUM (
    'PENDING',    -- In attesa verifica email
    'ACTIVE',     -- Attivo
    'SUSPENDED',  -- Sospeso temporaneamente
    'BANNED'      -- Bannato permanentemente
);

-- Tipo di annuncio
CREATE TYPE listing_type AS ENUM (
    'SALE',       -- Vendita normale
    'GIFT'        -- Regalo (gratis)
);

-- Metodo di consegna
CREATE TYPE shipping_method AS ENUM (
    'PICKUP',            -- Ritiro in loco
    'SELLER_SHIPS',      -- Il venditore spedisce
    'BUYER_ARRANGES',    -- Il compratore organizza il trasporto
    'PLATFORM_MANAGED'   -- GecoGreen gestisce la spedizione (futuro)
);

-- Stato prodotto
CREATE TYPE product_status AS ENUM (
    'DRAFT',      -- Bozza, non visibile
    'ACTIVE',     -- In vendita
    'SOLD',       -- Venduto
    'EXPIRED',    -- Scaduto (es. cibo)
    'DELETED'     -- Rimosso
);

-- Stato ordine
CREATE TYPE order_status AS ENUM (
    'CREATED',            -- Creato, in attesa pagamento
    'AWAITING_PAYMENT',   -- Checkout iniziato
    'PAID',               -- Pagato, in attesa ritiro/spedizione
    'READY_PICKUP',       -- Pronto per il ritiro
    'SHIPPED',            -- Spedito (tracking inserito)
    'DELIVERED',          -- Consegnato (da tracking o conferma)
    'COMPLETED',          -- Completato (soldi sbloccati)
    'DISPUTED',           -- Contestazione aperta
    'CANCELLED',          -- Annullato
    'REFUNDED'            -- Rimborsato
);

-- Motivo disputa
CREATE TYPE dispute_reason AS ENUM (
    'ITEM_NOT_RECEIVED',      -- Mai arrivato
    'ITEM_DAMAGED',           -- Danneggiato
    'ITEM_NOT_AS_DESCRIBED',  -- Diverso dalla descrizione
    'SELLER_NO_SHOW',         -- Venditore irreperibile
    'SCAM_ATTEMPT'            -- Tentativo di truffa
);

-- Stato disputa
CREATE TYPE dispute_status AS ENUM (
    'OPEN',                   -- Appena aperta
    'WAITING_SELLER',         -- In attesa risposta venditore
    'WAITING_BUYER',          -- In attesa risposta compratore
    'ADMIN_REVIEW',           -- In revisione admin
    'RESOLVED_REFUND',        -- Risolta con rimborso
    'RESOLVED_PAYOUT',        -- Risolta con pagamento al seller
    'RESOLVED_PARTIAL'        -- Risolta con rimborso parziale
);

-- Stato messaggio
CREATE TYPE message_status AS ENUM (
    'SENT',       -- Inviato
    'DELIVERED',  -- Consegnato
    'READ',       -- Letto
    'BLOCKED'     -- Bloccato dalla moderazione
);

-- ============================================
-- TABLES
-- ============================================

-- UTENTI
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,

    -- Dati personali
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone VARCHAR(20),

    -- Localizzazione
    city VARCHAR(100),
    province VARCHAR(50),
    postal_code VARCHAR(10),
    country VARCHAR(2) DEFAULT 'IT',

    -- Ruoli (un utente può avere più ruoli)
    roles user_role[] DEFAULT '{BUYER}',
    status user_status DEFAULT 'PENDING',

    -- Email verification
    email_verified BOOLEAN DEFAULT FALSE,
    email_verification_token VARCHAR(100),
    email_verification_expires TIMESTAMP,

    -- Password reset
    password_reset_token VARCHAR(100),
    password_reset_expires TIMESTAMP,

    -- OAuth / Social Login
    google_id VARCHAR(100) UNIQUE,         -- Google OAuth sub
    apple_id VARCHAR(100) UNIQUE,          -- Apple Sign In sub
    facebook_id VARCHAR(100) UNIQUE,       -- Facebook OAuth (futuro)
    oauth_provider VARCHAR(20),            -- 'google', 'apple', 'facebook', NULL se email
    avatar_url VARCHAR(500),               -- Foto profilo da OAuth

    -- Stripe
    stripe_customer_id VARCHAR(50),       -- Per i buyer
    stripe_account_id VARCHAR(50),        -- Per i seller (Connect)
    stripe_onboarding_complete BOOLEAN DEFAULT FALSE,

    -- Metriche
    rating_avg DECIMAL(3,2) DEFAULT 0,
    rating_count INTEGER DEFAULT 0,
    strike_count INTEGER DEFAULT 0,

    -- Impatto Ambientale (Gamification)
    total_co2_saved DECIMAL(10,2) DEFAULT 0,       -- kg CO2 risparmiati totali
    total_water_saved DECIMAL(10,2) DEFAULT 0,     -- litri acqua risparmiati
    eco_credits INT DEFAULT 0,                      -- Punti EcoCredits
    eco_level VARCHAR(50) DEFAULT 'Germoglio',      -- Livello gamification

    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login_at TIMESTAMP,

    -- Soft delete
    deleted_at TIMESTAMP
);

-- DATI AZIENDALI SELLER
CREATE TABLE seller_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,

    -- Dati aziendali
    business_name VARCHAR(255) NOT NULL,
    vat_number VARCHAR(20),                -- P.IVA
    fiscal_code VARCHAR(20),               -- Codice fiscale

    -- Indirizzo sede
    address_street VARCHAR(255),
    address_city VARCHAR(100),
    address_province VARCHAR(50),
    address_postal_code VARCHAR(10),
    address_country VARCHAR(2) DEFAULT 'IT',

    -- Coordinate per mappa
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),

    -- Info attività
    business_description TEXT,
    business_category VARCHAR(100),

    -- Orari ritiro
    pickup_hours JSONB,  -- {"mon": "9-18", "tue": "9-18", ...}

    -- Verifica
    verified BOOLEAN DEFAULT FALSE,
    verified_at TIMESTAMP,
    verified_by UUID REFERENCES users(id),

    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(user_id)
);

-- PUNTI VENDITA / SEDI RITIRO (multi-location per seller)
CREATE TABLE seller_locations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    seller_id UUID REFERENCES users(id) ON DELETE CASCADE,

    -- Info punto vendita
    name VARCHAR(255) NOT NULL,            -- Es: "Sede Milano Centro", "Magazzino Nord"
    is_primary BOOLEAN DEFAULT FALSE,      -- Sede principale
    is_active BOOLEAN DEFAULT TRUE,

    -- Indirizzo
    address_street VARCHAR(255) NOT NULL,
    address_city VARCHAR(100) NOT NULL,
    address_province VARCHAR(50),
    address_postal_code VARCHAR(10) NOT NULL,
    address_country VARCHAR(2) DEFAULT 'IT',

    -- Coordinate per mappa
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),

    -- Contatti specifici di questa sede
    phone VARCHAR(20),
    email VARCHAR(255),

    -- Orari ritiro (JSON per flessibilità)
    pickup_hours JSONB,  -- {"mon": "9-18", "tue": "9-18", "sat": "9-13", "sun": null}

    -- Note per il ritiro
    pickup_instructions TEXT,  -- Es: "Entrare dal retro", "Citofono 3"

    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indice per trovare tutte le sedi di un seller
CREATE INDEX idx_seller_locations_seller ON seller_locations(seller_id);
CREATE INDEX idx_seller_locations_city ON seller_locations(address_city);

-- CATEGORIE PRODOTTI
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    icon VARCHAR(50),
    parent_id UUID REFERENCES categories(id),
    sort_order INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,

    -- Valori impatto ambientale stimati (per calcolo certificati)
    estimated_co2_kg DECIMAL(10,2) DEFAULT 0,      -- kg CO2 risparmiati per unità
    estimated_water_l DECIMAL(10,2) DEFAULT 0,     -- litri acqua risparmiati per unità
    estimated_waste_kg DECIMAL(10,2) DEFAULT 0,    -- kg rifiuti evitati per unità

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- PRODOTTI
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    seller_id UUID REFERENCES users(id) ON DELETE CASCADE,
    category_id UUID REFERENCES categories(id),
    location_id UUID REFERENCES seller_locations(id),  -- Sede di ritiro specifica

    -- Info prodotto
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,

    -- Prezzi
    price DECIMAL(10, 2) NOT NULL,              -- Prezzo attuale
    original_price DECIMAL(10, 2),              -- Prezzo originale (per sconto)

    -- Tipo annuncio
    listing_type listing_type DEFAULT 'SALE',

    -- Spedizione
    shipping_method shipping_method DEFAULT 'PICKUP',
    shipping_cost DECIMAL(10, 2) DEFAULT 0,     -- Se SELLER_SHIPS

    -- Quantità
    quantity INTEGER DEFAULT 1,
    quantity_available INTEGER DEFAULT 1,

    -- Scadenza (per alimentari)
    expiry_date DATE,
    expiry_photo_url VARCHAR(500),        -- Foto specifica della data di scadenza

    -- Dutch Auction (prezzo che scende nel tempo)
    is_dutch_auction BOOLEAN DEFAULT FALSE,
    dutch_start_price DECIMAL(10, 2),     -- Prezzo iniziale (es: 100€)
    dutch_decrease_amount DECIMAL(10, 2), -- Quanto scende ogni intervallo (es: 1€)
    dutch_decrease_hours INT DEFAULT 24,  -- Intervallo in ore (default: 1 giorno)
    dutch_min_price DECIMAL(10, 2),       -- Prezzo minimo (stop decrescita)
    dutch_started_at TIMESTAMP,           -- Quando è partito il countdown

    -- Peso (per spedizioni)
    weight_kg DECIMAL(10, 2),

    -- Localizzazione (ereditata da seller o override)
    city VARCHAR(100),
    province VARCHAR(50),
    postal_code VARCHAR(10),
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),

    -- Immagini
    images JSONB DEFAULT '[]',  -- Array di URL

    -- Stato
    status product_status DEFAULT 'DRAFT',

    -- SEO
    slug VARCHAR(255),

    -- Metriche
    view_count INTEGER DEFAULT 0,
    favorite_count INTEGER DEFAULT 0,

    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    published_at TIMESTAMP,

    -- Soft delete
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

-- ORDINI
CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_number VARCHAR(20) UNIQUE NOT NULL,  -- Es: LZ-2024-000001

    -- Parti coinvolte
    buyer_id UUID REFERENCES users(id),
    seller_id UUID REFERENCES users(id),
    product_id UUID REFERENCES products(id),

    -- Quantità e prezzi
    quantity INTEGER NOT NULL DEFAULT 1,
    unit_price DECIMAL(10, 2) NOT NULL,
    shipping_cost DECIMAL(10, 2) DEFAULT 0,
    total_price DECIMAL(10, 2) NOT NULL,

    -- Commissione
    commission_rate DECIMAL(5, 4) DEFAULT 0.10,  -- 10%
    commission_amount DECIMAL(10, 2) NOT NULL,

    -- Metodo consegna
    shipping_method shipping_method NOT NULL,

    -- Indirizzo spedizione (se applicabile)
    shipping_address JSONB,

    -- Tracking (se spedito)
    tracking_number VARCHAR(100),
    tracking_url VARCHAR(500),

    -- Platform Managed Shipping (gestione corriere GecoGreen - futuro)
    platform_shipping_enabled BOOLEAN DEFAULT FALSE,  -- Se TRUE, GecoGreen gestisce
    platform_shipping_cost DECIMAL(10, 2),            -- Costo calcolato dalla piattaforma
    platform_courier_id VARCHAR(100),                 -- ID corriere (es: BRT, GLS, DHL)
    platform_courier_name VARCHAR(100),               -- Nome corriere
    platform_label_url VARCHAR(500),                  -- URL etichetta di spedizione
    platform_pickup_scheduled_at TIMESTAMP,           -- Data ritiro corriere dal seller

    -- QR Code per ritiro
    qr_code VARCHAR(100),                  -- QR code principale
    qr_code_text VARCHAR(20),              -- Codice testuale (ABC-123-XYZ)
    qr_code_expires TIMESTAMP,

    -- Delega ritiro
    delegate_name VARCHAR(200),            -- Nome completo delegato
    delegate_code VARCHAR(20),             -- Codice per il delegato (DEL-xxx)
    delegated_at TIMESTAMP,                -- Quando è stata creata la delega
    picked_up_by VARCHAR(50),              -- 'BUYER' o 'DELEGATE'
    picked_up_by_name VARCHAR(200),        -- Nome di chi ha effettivamente ritirato

    -- Scadenze
    deadline_date TIMESTAMP,  -- Deadline ritiro/conferma

    -- Stato
    status order_status DEFAULT 'CREATED',

    -- Pagamento
    stripe_payment_intent_id VARCHAR(100),
    stripe_transfer_id VARCHAR(100),
    paid_at TIMESTAMP,

    -- Payout
    payout_status VARCHAR(20) DEFAULT 'PENDING',
    payout_at TIMESTAMP,

    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP,
    cancelled_at TIMESTAMP,

    -- Note
    buyer_notes TEXT,
    seller_notes TEXT,

    -- Impatto Ambientale (calcolato al completamento ordine)
    co2_saved DECIMAL(10,2),               -- kg CO2 risparmiati
    water_saved DECIMAL(10,2),             -- litri acqua risparmiati
    waste_avoided DECIMAL(10,2),           -- kg rifiuti evitati
    eco_credits_earned INT DEFAULT 0,      -- EcoCredits guadagnati buyer
    impact_certificate_url VARCHAR(500)    -- URL del certificato PDF/PNG generato
);

-- DISPUTE
CREATE TABLE disputes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID REFERENCES orders(id) ON DELETE CASCADE,

    -- Chi ha aperto
    opened_by UUID REFERENCES users(id),
    opened_by_role user_role,

    -- Motivo
    reason dispute_reason NOT NULL,
    description TEXT NOT NULL,

    -- Prove
    evidence_photos JSONB DEFAULT '[]',  -- Array di URL

    -- Stato
    status dispute_status DEFAULT 'OPEN',

    -- Risoluzione
    resolution_type VARCHAR(50),
    resolution_amount DECIMAL(10, 2),  -- Se rimborso parziale
    resolution_notes TEXT,
    resolved_by UUID REFERENCES users(id),
    resolved_at TIMESTAMP,

    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(order_id)  -- Una sola disputa per ordine
);

-- MESSAGGI CHAT
CREATE TABLE messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID REFERENCES orders(id) ON DELETE CASCADE,

    -- Mittente
    sender_id UUID REFERENCES users(id),

    -- Contenuto
    content TEXT NOT NULL,

    -- Allegati
    attachments JSONB DEFAULT '[]',

    -- Stato
    status message_status DEFAULT 'SENT',

    -- Moderazione
    is_flagged BOOLEAN DEFAULT FALSE,
    flag_reason VARCHAR(100),

    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    read_at TIMESTAMP
);

-- FEEDBACK/RECENSIONI
CREATE TABLE reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID REFERENCES orders(id) ON DELETE CASCADE,

    -- Da chi a chi
    from_user_id UUID REFERENCES users(id),
    to_user_id UUID REFERENCES users(id),

    -- Valutazione
    rating INTEGER CHECK (rating >= 1 AND rating <= 5),
    comment TEXT,

    -- Visibilità
    is_anonymous BOOLEAN DEFAULT FALSE,
    is_visible BOOLEAN DEFAULT TRUE,

    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(order_id, from_user_id)  -- Un feedback per ordine per utente
);

-- COMMISSIONI (per fatturazione)
CREATE TABLE commissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID REFERENCES orders(id),
    seller_id UUID REFERENCES users(id),

    -- Importi
    gross_amount DECIMAL(10, 2) NOT NULL,      -- Prezzo vendita
    commission_rate DECIMAL(5, 4) NOT NULL,
    commission_amount DECIMAL(10, 2) NOT NULL,

    -- Fatturazione
    invoice_id UUID,  -- Riferimento a fattura mensile
    invoiced_at TIMESTAMP,

    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- FATTURE COMMISSIONI
CREATE TABLE invoices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_number VARCHAR(20) UNIQUE NOT NULL,  -- Es: INV-2024-000001
    seller_id UUID REFERENCES users(id),

    -- Periodo
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,

    -- Importi
    subtotal DECIMAL(10, 2) NOT NULL,
    vat_rate DECIMAL(5, 4) DEFAULT 0.22,
    vat_amount DECIMAL(10, 2) NOT NULL,
    total DECIMAL(10, 2) NOT NULL,

    -- PDF
    pdf_url VARCHAR(500),

    -- Stato
    status VARCHAR(20) DEFAULT 'DRAFT',  -- DRAFT, SENT, PAID
    sent_at TIMESTAMP,
    paid_at TIMESTAMP,

    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- AUDIT LOG
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Chi ha fatto l'azione
    user_id UUID REFERENCES users(id),

    -- Cosa
    action VARCHAR(100) NOT NULL,  -- Es: USER_BANNED, ORDER_REFUNDED
    entity_type VARCHAR(50),       -- Es: user, order, product
    entity_id UUID,

    -- Dettagli
    old_value JSONB,
    new_value JSONB,
    metadata JSONB,

    -- IP e device
    ip_address INET,
    user_agent TEXT,

    -- Timestamp
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- SESSIONI (per refresh token)
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,

    refresh_token VARCHAR(500) NOT NULL,
    expires_at TIMESTAMP NOT NULL,

    -- Device info
    device_info JSONB,
    ip_address INET,

    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_used_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- IMPACT LOGS (Storico azioni per EcoCredits e certificati)
CREATE TABLE impact_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    order_id UUID REFERENCES orders(id) ON DELETE SET NULL,

    -- Tipo azione
    action_type VARCHAR(50) NOT NULL,  -- 'PURCHASE', 'SALE', 'GIFT', 'BONUS', 'REDEEM', 'REFERRAL'

    -- Impatto ambientale
    co2_saved DECIMAL(10,2) DEFAULT 0,
    water_saved DECIMAL(10,2) DEFAULT 0,
    waste_avoided DECIMAL(10,2) DEFAULT 0,

    -- Punti EcoCredits
    eco_credits_earned INT DEFAULT 0,
    eco_credits_spent INT DEFAULT 0,

    -- Descrizione
    description TEXT,

    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- CSR REPORTS (Report sostenibilità per aziende)
CREATE TABLE csr_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    seller_id UUID REFERENCES users(id) ON DELETE CASCADE,

    -- Periodo di riferimento
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,

    -- Statistiche aggregate
    total_items INT DEFAULT 0,
    total_co2_kg DECIMAL(10,2) DEFAULT 0,
    total_water_l DECIMAL(10,2) DEFAULT 0,
    total_waste_kg DECIMAL(10,2) DEFAULT 0,
    total_value DECIMAL(10,2) DEFAULT 0,

    -- Dettaglio per categoria (JSON per flessibilità)
    category_breakdown JSONB,  -- [{"category": "Alimentari", "qty": 50, "co2": 100, "value": 500}, ...]

    -- File generato
    pdf_url VARCHAR(500),
    report_code VARCHAR(50) UNIQUE,  -- CSR-2024-AZIENDA-ABC123

    -- Piano (per pricing futuro)
    plan_type VARCHAR(20) DEFAULT 'BASE',  -- 'BASE', 'PRO', 'ENTERPRISE'

    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- PLATFORM SHIPPING CONFIG (Configurazione spedizioni gestite da GecoGreen - Futuro)
CREATE TABLE platform_shipping_config (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Info corriere
    courier_id VARCHAR(50) NOT NULL UNIQUE,    -- 'BRT', 'GLS', 'DHL', 'POSTE'
    courier_name VARCHAR(100) NOT NULL,
    courier_logo_url VARCHAR(500),

    -- Stato
    is_enabled BOOLEAN DEFAULT FALSE,          -- Admin abilita/disabilita

    -- Tariffe (semplificate - in futuro API corriere)
    base_rate DECIMAL(10, 2) NOT NULL,         -- Tariffa base (es: 5€)
    rate_per_kg DECIMAL(10, 2) DEFAULT 0,      -- Costo aggiuntivo per kg
    max_weight_kg DECIMAL(10, 2),              -- Peso massimo accettato

    -- Zone/Regioni abilitate (JSON per flessibilità)
    enabled_regions JSONB DEFAULT '["IT"]',    -- Es: ["IT"], ["IT", "FR", "DE"]

    -- Margine piattaforma (markup sul costo corriere)
    platform_markup_percent DECIMAL(5, 2) DEFAULT 10,  -- % markup GecoGreen

    -- Tempi di consegna stimati
    estimated_days_min INT DEFAULT 2,
    estimated_days_max INT DEFAULT 5,

    -- Note admin
    admin_notes TEXT,

    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- PLATFORM SHIPPING SETTINGS (Impostazioni globali admin)
CREATE TABLE platform_shipping_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Chiave-valore per flessibilità
    setting_key VARCHAR(100) NOT NULL UNIQUE,
    setting_value TEXT,
    setting_type VARCHAR(20) DEFAULT 'string',  -- 'string', 'boolean', 'number', 'json'

    -- Descrizione (per UI admin)
    description TEXT,

    -- Timestamps
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by UUID REFERENCES users(id)
);

-- Seed settings iniziali per shipping gestito
INSERT INTO platform_shipping_settings (setting_key, setting_value, setting_type, description) VALUES
('platform_shipping_enabled', 'false', 'boolean', 'Abilita/disabilita globalmente le spedizioni gestite da GecoGreen'),
('platform_shipping_min_order', '20', 'number', 'Ordine minimo per spedizione gestita (€)'),
('platform_shipping_free_above', '100', 'number', 'Spedizione gratuita sopra questo importo (€)'),
('platform_shipping_insurance_rate', '0.01', 'number', 'Percentuale assicurazione sul valore merce');

-- ECO REWARDS (Catalogo premi per spendere EcoCredits)
CREATE TABLE eco_rewards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Info premio
    name VARCHAR(200) NOT NULL,
    description TEXT,
    icon VARCHAR(50),

    -- Costo in EcoCredits
    eco_credits_cost INT NOT NULL,

    -- Tipo premio
    reward_type VARCHAR(50) NOT NULL,  -- 'COMMISSION_WAIVER', 'VISIBILITY_BOOST', 'BADGE', 'TREE_DONATION'

    -- Configurazione premio (JSON per flessibilità)
    config JSONB,  -- Es: {"duration_hours": 24} per boost, {"tree_count": 1} per donazione

    -- Disponibilità
    is_active BOOLEAN DEFAULT TRUE,
    stock INT,  -- NULL = illimitato

    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- REDEEMED REWARDS (Premi riscattati dagli utenti)
CREATE TABLE redeemed_rewards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    reward_id UUID REFERENCES eco_rewards(id),

    -- Punti spesi
    eco_credits_spent INT NOT NULL,

    -- Stato utilizzo
    status VARCHAR(20) DEFAULT 'PENDING',  -- 'PENDING', 'USED', 'EXPIRED'
    used_at TIMESTAMP,
    expires_at TIMESTAMP,

    -- Riferimento (es: ordine per cui è stato usato lo sconto commissioni)
    applied_to_order_id UUID REFERENCES orders(id),

    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- COMMUNITY IMPACT (Statistiche globali - cache aggregata)
CREATE TABLE community_impact (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Statistiche totali
    total_co2_saved DECIMAL(12,2) DEFAULT 0,
    total_water_saved DECIMAL(12,2) DEFAULT 0,
    total_waste_avoided DECIMAL(12,2) DEFAULT 0,
    total_orders_completed INT DEFAULT 0,
    total_items_saved INT DEFAULT 0,

    -- Alberi piantati (integrazione Tree-Nation)
    trees_planted INT DEFAULT 0,
    trees_target INT DEFAULT 0,  -- Prossimo obiettivo

    -- Ultimo aggiornamento
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================
-- INDEXES
-- ============================================

-- Users
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_stripe_customer ON users(stripe_customer_id);
CREATE INDEX idx_users_stripe_account ON users(stripe_account_id);
CREATE INDEX idx_users_google ON users(google_id);
CREATE INDEX idx_users_apple ON users(apple_id);

-- Seller profiles
CREATE INDEX idx_seller_profiles_user ON seller_profiles(user_id);
CREATE INDEX idx_seller_profiles_vat ON seller_profiles(vat_number);

-- Products
CREATE INDEX idx_products_seller ON products(seller_id);
CREATE INDEX idx_products_category ON products(category_id);
CREATE INDEX idx_products_location ON products(location_id);
CREATE INDEX idx_products_status ON products(status);
CREATE INDEX idx_products_listing_type ON products(listing_type);
CREATE INDEX idx_products_city ON products(city);
CREATE INDEX idx_products_expiry ON products(expiry_date);
CREATE INDEX idx_products_created ON products(created_at DESC);

-- Full text search su prodotti
CREATE INDEX idx_products_search ON products
    USING GIN (to_tsvector('italian', title || ' ' || description));

-- Orders
CREATE INDEX idx_orders_buyer ON orders(buyer_id);
CREATE INDEX idx_orders_seller ON orders(seller_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_created ON orders(created_at DESC);
CREATE INDEX idx_orders_number ON orders(order_number);

-- Messages
CREATE INDEX idx_messages_order ON messages(order_id);
CREATE INDEX idx_messages_sender ON messages(sender_id);
CREATE INDEX idx_messages_created ON messages(created_at);

-- Disputes
CREATE INDEX idx_disputes_order ON disputes(order_id);
CREATE INDEX idx_disputes_status ON disputes(status);

-- Reviews
CREATE INDEX idx_reviews_to_user ON reviews(to_user_id);
CREATE INDEX idx_reviews_from_user ON reviews(from_user_id);

-- Commissions
CREATE INDEX idx_commissions_seller ON commissions(seller_id);
CREATE INDEX idx_commissions_invoice ON commissions(invoice_id);

-- Sessions
CREATE INDEX idx_sessions_user ON sessions(user_id);
CREATE INDEX idx_sessions_token ON sessions(refresh_token);

-- Audit logs
CREATE INDEX idx_audit_user ON audit_logs(user_id);
CREATE INDEX idx_audit_entity ON audit_logs(entity_type, entity_id);
CREATE INDEX idx_audit_created ON audit_logs(created_at DESC);

-- Impact logs
CREATE INDEX idx_impact_logs_user ON impact_logs(user_id);
CREATE INDEX idx_impact_logs_order ON impact_logs(order_id);
CREATE INDEX idx_impact_logs_action ON impact_logs(action_type);
CREATE INDEX idx_impact_logs_created ON impact_logs(created_at DESC);

-- CSR Reports
CREATE INDEX idx_csr_reports_seller ON csr_reports(seller_id);
CREATE INDEX idx_csr_reports_period ON csr_reports(period_start, period_end);
CREATE INDEX idx_csr_reports_code ON csr_reports(report_code);

-- Eco Rewards
CREATE INDEX idx_eco_rewards_type ON eco_rewards(reward_type);
CREATE INDEX idx_eco_rewards_active ON eco_rewards(is_active);

-- Redeemed Rewards
CREATE INDEX idx_redeemed_rewards_user ON redeemed_rewards(user_id);
CREATE INDEX idx_redeemed_rewards_status ON redeemed_rewards(status);

-- Platform Shipping Config
CREATE INDEX idx_platform_shipping_config_enabled ON platform_shipping_config(is_enabled);
CREATE INDEX idx_platform_shipping_config_courier ON platform_shipping_config(courier_id);

-- Products - Dutch Auction
CREATE INDEX idx_products_dutch_auction ON products(is_dutch_auction) WHERE is_dutch_auction = TRUE;

-- ============================================
-- TRIGGERS
-- ============================================

-- Funzione per aggiornare updated_at
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Applica trigger a tutte le tabelle con updated_at
CREATE TRIGGER users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER seller_profiles_updated_at
    BEFORE UPDATE ON seller_profiles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER products_updated_at
    BEFORE UPDATE ON products
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER orders_updated_at
    BEFORE UPDATE ON orders
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER disputes_updated_at
    BEFORE UPDATE ON disputes
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER reviews_updated_at
    BEFORE UPDATE ON reviews
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- ============================================
-- SEED DATA (Categorie e Sottocategorie)
-- ============================================

-- CATEGORIE PRINCIPALI (14)
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

-- SOTTOCATEGORIE: Alimentari Freschi
INSERT INTO categories (name, slug, parent_id, sort_order) VALUES
('Latticini', 'latticini', '10000000-0000-0000-0000-000000000001', 1),
('Latte e Panna', 'latte-panna', '10000000-0000-0000-0000-000000000001', 2),
('Carne Bovina', 'carne-bovina', '10000000-0000-0000-0000-000000000001', 3),
('Carne Suina', 'carne-suina', '10000000-0000-0000-0000-000000000001', 4),
('Pollame', 'pollame', '10000000-0000-0000-0000-000000000001', 5),
('Pesce Fresco', 'pesce-fresco', '10000000-0000-0000-0000-000000000001', 6),
('Uova', 'uova', '10000000-0000-0000-0000-000000000001', 7),
('Salumi Freschi', 'salumi-freschi', '10000000-0000-0000-0000-000000000001', 8),
('Pasta Fresca', 'pasta-fresca', '10000000-0000-0000-0000-000000000001', 9),
('Gastronomia', 'gastronomia', '10000000-0000-0000-0000-000000000001', 10),
('Pane e Pizze Fresche', 'pane-pizze', '10000000-0000-0000-0000-000000000001', 11);

-- SOTTOCATEGORIE: Alimentari Confezionati
INSERT INTO categories (name, slug, parent_id, sort_order) VALUES
('Conserve e Scatolame', 'conserve', '10000000-0000-0000-0000-000000000002', 1),
('Pasta Secca', 'pasta-secca', '10000000-0000-0000-0000-000000000002', 2),
('Riso e Cereali', 'riso-cereali', '10000000-0000-0000-0000-000000000002', 3),
('Biscotti e Snack', 'biscotti-snack', '10000000-0000-0000-0000-000000000002', 4),
('Salse e Condimenti', 'salse-condimenti', '10000000-0000-0000-0000-000000000002', 5),
('Farine e Lieviti', 'farine-lieviti', '10000000-0000-0000-0000-000000000002', 6),
('Caffè e Tè', 'caffe-te', '10000000-0000-0000-0000-000000000002', 7),
('Cioccolato e Dolci', 'cioccolato-dolci', '10000000-0000-0000-0000-000000000002', 8),
('Olio e Aceto', 'olio-aceto', '10000000-0000-0000-0000-000000000002', 9),
('Spezie', 'spezie', '10000000-0000-0000-0000-000000000002', 10),
('Affettati Confezionati', 'affettati-confezionati', '10000000-0000-0000-0000-000000000002', 11);

-- SOTTOCATEGORIE: Bevande
INSERT INTO categories (name, slug, parent_id, sort_order) VALUES
('Latte UHT', 'latte-uht', '10000000-0000-0000-0000-000000000003', 1),
('Bevande Vegetali', 'bevande-vegetali', '10000000-0000-0000-0000-000000000003', 2),
('Succhi di Frutta', 'succhi-frutta', '10000000-0000-0000-0000-000000000003', 3),
('Bibite Gassate', 'bibite-gassate', '10000000-0000-0000-0000-000000000003', 4),
('Birra', 'birra', '10000000-0000-0000-0000-000000000003', 5),
('Acqua', 'acqua', '10000000-0000-0000-0000-000000000003', 6),
('Energy Drink', 'energy-drink', '10000000-0000-0000-0000-000000000003', 7),
('Integratori Liquidi', 'integratori-liquidi', '10000000-0000-0000-0000-000000000003', 8);

-- SOTTOCATEGORIE: Frutta e Verdura
INSERT INTO categories (name, slug, parent_id, sort_order) VALUES
('Frutta Fresca', 'frutta-fresca', '10000000-0000-0000-0000-000000000004', 1),
('Verdura Fresca', 'verdura-fresca', '10000000-0000-0000-0000-000000000004', 2),
('Quarta Gamma', 'quarta-gamma', '10000000-0000-0000-0000-000000000004', 3),
('Frutta Secca', 'frutta-secca', '10000000-0000-0000-0000-000000000004', 4);

-- SOTTOCATEGORIE: Surgelati
INSERT INTO categories (name, slug, parent_id, sort_order) VALUES
('Pesce Surgelato', 'pesce-surgelato', '10000000-0000-0000-0000-000000000005', 1),
('Carne Surgelata', 'carne-surgelata', '10000000-0000-0000-0000-000000000005', 2),
('Verdure Surgelate', 'verdure-surgelate', '10000000-0000-0000-0000-000000000005', 3),
('Pizze Surgelate', 'pizze-surgelate', '10000000-0000-0000-0000-000000000005', 4),
('Piatti Pronti', 'piatti-pronti', '10000000-0000-0000-0000-000000000005', 5),
('Gelati', 'gelati', '10000000-0000-0000-0000-000000000005', 6);

-- SOTTOCATEGORIE: Cosmetici
INSERT INTO categories (name, slug, parent_id, sort_order) VALUES
('Creme Viso', 'creme-viso', '10000000-0000-0000-0000-000000000006', 1),
('Creme Corpo', 'creme-corpo', '10000000-0000-0000-0000-000000000006', 2),
('Trucco Viso', 'trucco-viso', '10000000-0000-0000-0000-000000000006', 3),
('Trucco Occhi', 'trucco-occhi', '10000000-0000-0000-0000-000000000006', 4),
('Trucco Labbra', 'trucco-labbra', '10000000-0000-0000-0000-000000000006', 5),
('Solari', 'solari', '10000000-0000-0000-0000-000000000006', 6),
('Profumi', 'profumi', '10000000-0000-0000-0000-000000000006', 7);

-- SOTTOCATEGORIE: Cura Persona
INSERT INTO categories (name, slug, parent_id, sort_order) VALUES
('Shampoo e Balsamo', 'shampoo-balsamo', '10000000-0000-0000-0000-000000000007', 1),
('Bagnoschiuma', 'bagnoschiuma', '10000000-0000-0000-0000-000000000007', 2),
('Dentifrici', 'dentifrici', '10000000-0000-0000-0000-000000000007', 3),
('Collutori', 'collutori', '10000000-0000-0000-0000-000000000007', 4),
('Deodoranti', 'deodoranti', '10000000-0000-0000-0000-000000000007', 5),
('Rasatura', 'rasatura', '10000000-0000-0000-0000-000000000007', 6);

-- SOTTOCATEGORIE: Detergenza Casa
INSERT INTO categories (name, slug, parent_id, sort_order) VALUES
('Detersivi Lavatrice', 'detersivi-lavatrice', '10000000-0000-0000-0000-000000000008', 1),
('Detersivi Piatti', 'detersivi-piatti', '10000000-0000-0000-0000-000000000008', 2),
('Candeggina', 'candeggina', '10000000-0000-0000-0000-000000000008', 3),
('Anticalcare', 'anticalcare', '10000000-0000-0000-0000-000000000008', 4),
('Disinfettanti', 'disinfettanti', '10000000-0000-0000-0000-000000000008', 5),
('Pavimenti', 'pavimenti', '10000000-0000-0000-0000-000000000008', 6);

-- SOTTOCATEGORIE: Pet Food
INSERT INTO categories (name, slug, parent_id, sort_order) VALUES
('Cibo Cani Umido', 'cibo-cani-umido', '10000000-0000-0000-0000-000000000009', 1),
('Cibo Cani Secco', 'cibo-cani-secco', '10000000-0000-0000-0000-000000000009', 2),
('Cibo Gatti Umido', 'cibo-gatti-umido', '10000000-0000-0000-0000-000000000009', 3),
('Cibo Gatti Secco', 'cibo-gatti-secco', '10000000-0000-0000-0000-000000000009', 4),
('Snack Animali', 'snack-animali', '10000000-0000-0000-0000-000000000009', 5),
('Altri Animali', 'altri-animali', '10000000-0000-0000-0000-000000000009', 6);

-- SOTTOCATEGORIE: Giardinaggio
INSERT INTO categories (name, slug, parent_id, sort_order) VALUES
('Fertilizzanti', 'fertilizzanti', '10000000-0000-0000-0000-000000000010', 1),
('Semi', 'semi', '10000000-0000-0000-0000-000000000010', 2),
('Fitofarmaci', 'fitofarmaci', '10000000-0000-0000-0000-000000000010', 3),
('Terricci', 'terricci', '10000000-0000-0000-0000-000000000010', 4),
('Diserbanti', 'diserbanti', '10000000-0000-0000-0000-000000000010', 5);

-- SOTTOCATEGORIE: Materiali Tecnici
INSERT INTO categories (name, slug, parent_id, sort_order) VALUES
('Vernici', 'vernici', '10000000-0000-0000-0000-000000000011', 1),
('Colle e Adesivi', 'colle-adesivi', '10000000-0000-0000-0000-000000000011', 2),
('Siliconi e Sigillanti', 'siliconi-sigillanti', '10000000-0000-0000-0000-000000000011', 3),
('Resine', 'resine', '10000000-0000-0000-0000-000000000011', 4),
('Stucchi', 'stucchi', '10000000-0000-0000-0000-000000000011', 5);

-- SOTTOCATEGORIE: Automotive
INSERT INTO categories (name, slug, parent_id, sort_order) VALUES
('Pneumatici', 'pneumatici', '10000000-0000-0000-0000-000000000012', 1),
('Olio Motore', 'olio-motore', '10000000-0000-0000-0000-000000000012', 2),
('Liquidi Auto', 'liquidi-auto', '10000000-0000-0000-0000-000000000012', 3),
('Batterie', 'batterie', '10000000-0000-0000-0000-000000000012', 4),
('Accessori Auto', 'accessori-auto', '10000000-0000-0000-0000-000000000012', 5);

-- SOTTOCATEGORIE: Sicurezza e DPI
INSERT INTO categories (name, slug, parent_id, sort_order) VALUES
('Estintori', 'estintori', '10000000-0000-0000-0000-000000000013', 1),
('Caschi', 'caschi', '10000000-0000-0000-0000-000000000013', 2),
('Imbragature', 'imbragature', '10000000-0000-0000-0000-000000000013', 3),
('Kit Pronto Soccorso', 'kit-pronto-soccorso', '10000000-0000-0000-0000-000000000013', 4),
('Guanti e Protezioni', 'guanti-protezioni', '10000000-0000-0000-0000-000000000013', 5);

-- SOTTOCATEGORIE: HORECA
INSERT INTO categories (name, slug, parent_id, sort_order) VALUES
('Ingredienti Ristorazione', 'ingredienti-ristorazione', '10000000-0000-0000-0000-000000000014', 1),
('Monoporzioni', 'monoporzioni', '10000000-0000-0000-0000-000000000014', 2),
('Preparati Gastronomici', 'preparati-gastronomici', '10000000-0000-0000-0000-000000000014', 3),
('Salse e Topping', 'salse-topping', '10000000-0000-0000-0000-000000000014', 4),
('Pasticceria', 'pasticceria', '10000000-0000-0000-0000-000000000014', 5);

-- AGGIORNAMENTO VALORI IMPATTO PER CATEGORIE PRINCIPALI
-- Valori medi stimati per unità di prodotto salvata
UPDATE categories SET estimated_co2_kg = 2, estimated_water_l = 500, estimated_waste_kg = 1
WHERE id = '10000000-0000-0000-0000-000000000001';  -- Alimentari Freschi

UPDATE categories SET estimated_co2_kg = 1, estimated_water_l = 200, estimated_waste_kg = 0.5
WHERE id = '10000000-0000-0000-0000-000000000002';  -- Alimentari Confezionati

UPDATE categories SET estimated_co2_kg = 0.5, estimated_water_l = 100, estimated_waste_kg = 0.3
WHERE id = '10000000-0000-0000-0000-000000000003';  -- Bevande

UPDATE categories SET estimated_co2_kg = 1, estimated_water_l = 300, estimated_waste_kg = 0.5
WHERE id = '10000000-0000-0000-0000-000000000004';  -- Frutta e Verdura

UPDATE categories SET estimated_co2_kg = 3, estimated_water_l = 400, estimated_waste_kg = 1
WHERE id = '10000000-0000-0000-0000-000000000005';  -- Surgelati

UPDATE categories SET estimated_co2_kg = 5, estimated_water_l = 1000, estimated_waste_kg = 0.5
WHERE id = '10000000-0000-0000-0000-000000000006';  -- Cosmetici

UPDATE categories SET estimated_co2_kg = 3, estimated_water_l = 800, estimated_waste_kg = 0.3
WHERE id = '10000000-0000-0000-0000-000000000007';  -- Cura Persona

UPDATE categories SET estimated_co2_kg = 4, estimated_water_l = 600, estimated_waste_kg = 1
WHERE id = '10000000-0000-0000-0000-000000000008';  -- Detergenza Casa

UPDATE categories SET estimated_co2_kg = 2, estimated_water_l = 400, estimated_waste_kg = 0.5
WHERE id = '10000000-0000-0000-0000-000000000009';  -- Pet Food

UPDATE categories SET estimated_co2_kg = 5, estimated_water_l = 200, estimated_waste_kg = 2
WHERE id = '10000000-0000-0000-0000-000000000010';  -- Giardinaggio

UPDATE categories SET estimated_co2_kg = 15, estimated_water_l = 500, estimated_waste_kg = 5
WHERE id = '10000000-0000-0000-0000-000000000011';  -- Materiali Tecnici

UPDATE categories SET estimated_co2_kg = 30, estimated_water_l = 2000, estimated_waste_kg = 10
WHERE id = '10000000-0000-0000-0000-000000000012';  -- Automotive

UPDATE categories SET estimated_co2_kg = 10, estimated_water_l = 800, estimated_waste_kg = 3
WHERE id = '10000000-0000-0000-0000-000000000013';  -- Sicurezza e DPI

UPDATE categories SET estimated_co2_kg = 3, estimated_water_l = 500, estimated_waste_kg = 1
WHERE id = '10000000-0000-0000-0000-000000000014';  -- HORECA

-- SEED DATA: Premi EcoCredits (solo premi economicamente sostenibili)
INSERT INTO eco_rewards (name, description, icon, eco_credits_cost, reward_type, config) VALUES
('Boost Visibilità 24h', 'Il tuo prodotto in evidenza per 24 ore', '🚀', 200, 'VISIBILITY_BOOST', '{"duration_hours": 24}'),
('Pianta un Albero', 'Donazione a Tree-Nation per piantare 1 albero', '🌳', 300, 'TREE_DONATION', '{"tree_count": 1, "partner": "tree-nation"}'),
('Badge Personalizzato', 'Un badge esclusivo sul tuo profilo', '🏅', 1000, 'BADGE', '{"badge_name": "Eco-Master"}');

-- SEED DATA: Record iniziale Community Impact
INSERT INTO community_impact (total_co2_saved, total_water_saved, total_waste_avoided, total_orders_completed, total_items_saved, trees_planted, trees_target)
VALUES (0, 0, 0, 0, 0, 0, 100);

-- ============================================
-- NOTES
-- ============================================

-- Per eseguire questo script:
-- psql -U gecogreen -d gecogreen -f database.sql

-- Per reset completo (ATTENZIONE: cancella tutto!):
-- DROP SCHEMA public CASCADE;
-- CREATE SCHEMA public;
-- GRANT ALL ON SCHEMA public TO gecogreen;
-- GRANT ALL ON SCHEMA public TO public;
-- \i database.sql

-- Per backup:
-- pg_dump -U gecogreen gecogreen > backup_$(date +%Y%m%d).sql

-- Per restore:
-- psql -U gecogreen gecogreen < backup_YYYYMMDD.sql
