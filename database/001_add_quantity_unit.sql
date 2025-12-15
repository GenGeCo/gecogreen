-- ============================================
-- Migration: Add quantity unit to products
-- ============================================

-- Add unit type enum
CREATE TYPE quantity_unit AS ENUM ('PIECE', 'KG', 'G', 'L', 'ML', 'CUSTOM');

-- Add columns to products table
ALTER TABLE products
ADD COLUMN quantity_unit quantity_unit DEFAULT 'PIECE',
ADD COLUMN quantity_unit_custom VARCHAR(50);

-- Update existing products to use PIECE as default
UPDATE products SET quantity_unit = 'PIECE' WHERE quantity_unit IS NULL;

COMMENT ON COLUMN products.quantity_unit IS 'Unità di misura per la quantità';
COMMENT ON COLUMN products.quantity_unit_custom IS 'Unità personalizzata quando quantity_unit = CUSTOM';
