-- ============================================
-- Migration: Add billing information to users
-- ============================================

-- Add billing fields to users table
ALTER TABLE users
ADD COLUMN fiscal_code VARCHAR(16),           -- Codice Fiscale (Italia)
ADD COLUMN sdi_code VARCHAR(7),               -- Codice Univoco SDI (fattura elettronica)
ADD COLUMN pec_email VARCHAR(255),            -- PEC per fattura elettronica
ADD COLUMN eu_vat_id VARCHAR(20),             -- VAT ID europeo
ADD COLUMN billing_address TEXT,              -- Indirizzo fatturazione
ADD COLUMN billing_city VARCHAR(100),         -- Città fatturazione
ADD COLUMN billing_province VARCHAR(2),       -- Provincia fatturazione
ADD COLUMN billing_postal_code VARCHAR(10),   -- CAP fatturazione
ADD COLUMN billing_country VARCHAR(2) DEFAULT 'IT'; -- Paese fatturazione (ISO 3166-1 alpha-2)

COMMENT ON COLUMN users.fiscal_code IS 'Codice Fiscale italiano';
COMMENT ON COLUMN users.sdi_code IS 'Codice Univoco SDI per fatturazione elettronica';
COMMENT ON COLUMN users.pec_email IS 'PEC per fatturazione elettronica (alternativa a SDI)';
COMMENT ON COLUMN users.eu_vat_id IS 'VAT ID europeo per clienti UE';
COMMENT ON COLUMN users.billing_address IS 'Indirizzo completo per fatturazione';
COMMENT ON COLUMN users.billing_city IS 'Città per fatturazione';
COMMENT ON COLUMN users.billing_province IS 'Provincia per fatturazione';
COMMENT ON COLUMN users.billing_postal_code IS 'CAP per fatturazione';
COMMENT ON COLUMN users.billing_country IS 'Paese per fatturazione (ISO 3166-1 alpha-2)';
